# 动态Go链接器

作用是能够动态的链接go编译后的.o后缀的二进制文件，配合热更新会有奇妙的效果。

**测试环境Go 1.10成功**

### 测试流程

- 首先在 $GOROOT/src/cmd文件夹下创建空文件夹 objfile。
- 然后将 $GOROOT/src/cmd/internal 文件夹拷贝到 $GOROOT/src/cmd/objfile 中。
- 若未执行上面两步，将编译main.go会报错。
- 然后进入bin中，修改测试用的request.go文件。
- 然后执行 编译.bat，得到 .o 后缀的二进制文件。
- 最后执行main.exe，它会动态加载 .o 的二进制文件。