#### 有一个文本1.txt, 每行结构为 "姓名 手机号 地址"， 用linux命令去重输出手机号
```shell
awk {print $2} 1.txt | uniq
```