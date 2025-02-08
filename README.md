# AICodeScan

![AICodeScan](https://socialify.git.ci/zacarx/AICodeScan/image?description=1&font=Source+Code+Pro&language=1&name=1&owner=1&pattern=Overlapping+Hexagons&theme=Dark)



## 工具概述

​	该工具基于Zjackky/CodeScan开发，通过对大多数不完整的代码以及依赖快速进行Sink点匹配，并且由AI进行审计精准定位，来帮助红队完成快速代码审计，目前工具支持的语言有PHP，Java。



## 编译

```bash
./build.sh
#需要golang环境
# 会生成所有版本在releases下
```



## 功能

​	在Zjackky/CodeScan基础上，增加ai审计。

​	sink点查找相关细节查看：https://github.com/xiaoqiuxx/CodeScan_public



## 使用

在程序当前目录增加config.yaml

config.yaml内容

```yaml
api:
  #  api接口base url
  url: "https://api.siliconflow.cn/v1/chat/completions"
  #  密钥，最好参考api文档，以下是轨迹流动的api
  key: "Bearer sk-key"

settings:
  #  每次调用ai间隔时间，防止频繁或者封号
  sleep_seconds: 3

model:
  #  模型名称
  name: "deepseek-ai/DeepSeek-R1-Distill-Qwen-14B"

#  这里%s不要动，防止输入错误
prompt:
  text: "请分析以下代码是否存在安全问题：\n文件: %s\n行号: %d\n内容:\n%s\n当前行：%s，请简明扼要，如果觉得大概率没有漏洞直接回答大概率没有漏洞七个汉字，如果有，严格按照一下格式输出：\n漏洞类型：\n危害等级：\n判断理由：\n可能的payload:"
```

下面是命令行用法

```bash
Usage of ./AICodeScan:
  -L string
        审计语言
  -d string
        要扫描的目录
  -h string
        使用帮助
  -lb string
        行黑名单
  -m string
        过滤的字符串
  -pb string
        路径黑名单
  -r string
        RCE规则
  -u string
        文件上传规则


Example:
	AICodeScan -L java -d ./net
	AICodeScan -L php -d ./net
	AICodeScan -d ./net -m "CheckSession.jsp"
```



## 效果图

![image-20250208130256764](./img/test.png)

## 联系

如果你想反馈体验，有更好的建议，或者想参与项目开发

可以添加我的微信（图一）

如果你想关注项目更新可以关注微信公众号（图二）

<img src="http://img.zacarx.top/img/image-20250208131538819.png" alt="image-20250208131538819" style="zoom:25%;" /><img src="http://img.zacarx.top/img/image-20250208131442284.png" alt="image-20250208131442284" style="zoom:50%;" />



## 声明

禁止使用本工具从事任何违法行为，使用本工具造成的一切后果应由使用工具当事人承担。

## 鸣谢

[xiaoqiuxx(github.com)](https://github.com/xiaoqiuxx)
