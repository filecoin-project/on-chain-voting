# I.Compiling Ucan Tool



## 1.First, you need to install the Go toolchain. You can find [instructions](https://go.dev/doc/install) here, with Go version >= 1.20



## 2. Get the code of UCAN signature tool

```bash
git clone https://github.com/black-domain/ucan-utils.git
```

## 3.Install dependencies

```bash
go mod tidy
```

## 4.Build the binary file

```bash
go build -o signature .
```

## 5.Now you should be able to run signature

```bash
$ ./signature -h
```

<img src="./img/1.png" style="zoom:50%;" />
