# Huawei bruteforce unlock bootloader

Note: This is a minimal go implementation of the python script located [here](https://github.com/SkyEmie/huawei-honor-unlock-bootloader).  

## requirements (build)

* golang >= 1.12

## build

```bash
go mod vendor -v
go build -ldflags '-s -w' -mod vendor -v -i -o build/unlock cmd/unlock/unlock.go
``` 

## requirements (run)

* adb
* fastboot
* devices booted to bootloader (`adb reboot bootloader`)

## run

```bash
./build/unlock --imei=123456789012345 --resume
``` 

The `--resume` option will create a resume file in case your script crashed or you stopped it. Of course, it will read from this file if you want to continue.

The `--oem-code` option can be used if you want to start on a specific oem code. Note that the resume file overwrites this value!
 
