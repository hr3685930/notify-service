package utils

import (
    "bytes"
    "errors"
    "io"
    "mime/multipart"
    "os"
    "os/exec"
    "path/filepath"
    "reflect"
    "runtime"
    "strconv"
    "strings"
    "sync"
)

// 当前项目根目录
var API_ROOT string

// 获取项目路径
func GetPath() string {

    if API_ROOT != "" {
        return API_ROOT
    }

    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        print(err.Error())
    }

    API_ROOT = strings.Replace(dir, "\\", "/", -1)
    return API_ROOT
}

// 判断文件目录否存在
func IsDirExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return fi.IsDir()
    }

}

// 创建文件
func MkdirFile(path string) error {

    err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
    if err != nil {
        return err
    }
    return nil
}

// DisableCache will disable caching of the home directory. Caching is enabled
// by default.
var DisableCache bool

var homedirCache string
var cacheLock sync.RWMutex

// Dir returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Dir() (string, error) {
    if !DisableCache {
        cacheLock.RLock()
        cached := homedirCache
        cacheLock.RUnlock()
        if cached != "" {
            return cached, nil
        }
    }

    cacheLock.Lock()
    defer cacheLock.Unlock()

    var result string
    var err error
    if runtime.GOOS == "windows" {
        result, err = dirWindows()
    } else {
        // Unix-like system, so just assume Unix
        result, err = dirUnix()
    }

    if err != nil {
        return "", err
    }
    homedirCache = result
    return result, nil
}

// Expand expands the path to include the home directory if the path
// is prefixed with `~`. If it isn't prefixed with `~`, the path is
// returned as-is.
func Expand(path string) (string, error) {
    if len(path) == 0 {
        return path, nil
    }

    if path[0] != '~' {
        return path, nil
    }

    if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
        return "", errors.New("cannot expand user-specific home dir")
    }

    dir, err := Dir()
    if err != nil {
        return "", err
    }

    return filepath.Join(dir, path[1:]), nil
}

func dirUnix() (string, error) {
    homeEnv := "HOME"
    if runtime.GOOS == "plan9" {
        // On plan9, env vars are lowercase.
        homeEnv = "home"
    }

    // First prefer the HOME environmental variable
    if home := os.Getenv(homeEnv); home != "" {
        return home, nil
    }

    var stdout bytes.Buffer

    // If that fails, try OS specific commands
    if runtime.GOOS == "darwin" {
        cmd := exec.Command("sh", "-c", `dscl -q . -read /Users/"$(whoami)" NFSHomeDirectory | sed 's/^[^ ]*: //'`)
        cmd.Stdout = &stdout
        if err := cmd.Run(); err == nil {
            result := strings.TrimSpace(stdout.String())
            if result != "" {
                return result, nil
            }
        }
    } else {
        cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
        cmd.Stdout = &stdout
        if err := cmd.Run(); err != nil {
            // If the error is ErrNotFound, we ignore it. Otherwise, return it.
            if err != exec.ErrNotFound {
                return "", err
            }
        } else {
            if passwd := strings.TrimSpace(stdout.String()); passwd != "" {
                // username:password:uid:gid:gecos:home:shell
                passwdParts := strings.SplitN(passwd, ":", 7)
                if len(passwdParts) > 5 {
                    return passwdParts[5], nil
                }
            }
        }
    }

    // If all else fails, try the shell
    stdout.Reset()
    cmd := exec.Command("sh", "-c", "cd && pwd")
    cmd.Stdout = &stdout
    if err := cmd.Run(); err != nil {
        return "", err
    }

    result := strings.TrimSpace(stdout.String())
    if result == "" {
        return "", errors.New("blank output when reading home directory")
    }

    return result, nil
}

func dirWindows() (string, error) {
    // First prefer the HOME environmental variable
    if home := os.Getenv("HOME"); home != "" {
        return home, nil
    }

    drive := os.Getenv("HOMEDRIVE")
    path := os.Getenv("HOMEPATH")
    home := drive + path
    if drive == "" || path == "" {
        home = os.Getenv("USERPROFILE")
    }
    if home == "" {
        return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
    }

    return home, nil
}

func DeleteFile(filePath string) error {
    return os.RemoveAll(filePath)
}

//创建文件夹,支持x/a/a  多层级
func MkDir(path string) error {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            //文件夹不存在，创建
            err = os.MkdirAll(path, os.ModePerm)
            if err != nil {
                return err
            }
        } else {
            return err
        }
    }
    return nil
}

//上传文件
func UploadFile(file *multipart.FileHeader, path string) (string, error) {
    if reflect.ValueOf(file).IsNil() || !reflect.ValueOf(file).IsValid() {
        return "", errors.New("invalid memory address or nil pointer dereference")
    }
    src, err := file.Open()
    defer src.Close()
    if err != nil {
        return "", err
    }
    err = MkDir(path)
    if err != nil {
        return "", err
    }
    // Destination
    // 去除空格
    filename := strings.Replace(file.Filename, " ", "", -1)
    // 去除换行符
    filename = strings.Replace(filename, "\n", "", -1)

    dst, err := os.Create(path + filename)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    // Copy
    if _, err = io.Copy(dst, src); err != nil {
        return "", err
    }
    return filename, nil
}

// 文件大小
func GetFileSize(filePath string) (int64, error) {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        return 0, err
    }
    fsize := fileInfo.Size()
    return fsize, nil
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
    var exist = true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false
    }
    return exist
}
