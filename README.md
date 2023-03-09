# ChatGPT-API

---

这是一个基于gin框架构建的chatGPT API服务

支持的功能有：
* 模型参数配置
* 提问上下文
* 非流式传输
* 流式传输

## 开始
### 下载源码
```shell
git clone https://github.com/xs818/openai-chatGPT-server.git
```

### 源码编译运行
```shell

# 复制配置文件，并修改自己的API KEY
cp config/example_config.toml config/dev_config.toml

go mod tidy

# 运行
go run main.go

```

### Docker运行
**镜像编译**
```shell
docker build -t openai-server:latest .
```

**linux系统运行**
```shell
# linux系统
docker run -d \
      --name openai-server \
      -p 8090:8090 \
      -v `pwd`/logs:/app/logs/ \
      -e CHAT_PORT=8090 \
      -e CHAT_APIKEY="your api_key" \
      openai-server:latest


```
`提示：` 
* 指定环境变量需要以`CHAT_`为前缀，并且环境变量需大写 
* 优先使用环境变量的配置


**linux系统下使用代理运行**
```shell
docker run --net=host -d \
          --name openai-server \
          -v `pwd`/logs:/app/logs/ \
          -e CHAT_PORT=8091 \
          -e CHAT_PROXY=127.0.0.1:7890 \
          openai-server:latest
```
`提示`： `CHAT_PROXY`需改成自己的代理地址。

**其他系统使用代理运行**
```shell
docker run -d \
            --name openai-server \
            -p 8090:8090 \
            -v `pwd`/logs:/app/logs/ \
            -e CHAT_PORT=8090 \
            -e CHAT_PROXY=http://host.docker.internal:7890 \
            --add-host=host.docker.internal:host-gateway \
            openai-server:latest
```

`提示`： CHAT_PROXY中的`http://host.docker.internal`是指向宿主机的IP,端口需要改成自己的代理端口

---
## 参数解释
#### Temperature
Temperature参数在ChatGPT中控制了生成文本的多样性。具体而言，它可以控制生成文本的随机性和创造性。较低的温度值会导致生成的文本较为保守，而较高的温度值则会导致生成的文本更加随机和创新。因此，温度值可以用来控制生成文本的风格和多样性。

在ChatGPT中，温度值通常在0.5到1.5之间取值，其中1.0是最常用的默认值。当温度值为1.0时，生成的文本具有一定的多样性，同时也能保持一定的连贯性和合理性。如果想要生成更加创新和出人意料的文本，可以将温度值调高；如果想要生成更加保守和符合预期的文本，可以将温度值调低。

需要注意的是，过高的温度值可能会导致生成的文本质量下降，因为生成的文本可能会变得不连贯、不合理甚至无意义。因此，在使用ChatGPT生成文本时，需要根据具体情况灵活调整温度值以达到最佳效果。


### top_p
top_p（也称为nucleus sampling或基于文本长度的筛选法），它可以用来控制生成文本的多样性和可控性。

具体来说，top_p参数是一个浮点数，它指定了从生成的概率分布中选择的概率质量上限。在生成文本时，ChatGPT会对每个词计算其在当前上下文下的条件概率，然后根据top_p参数筛选出概率质量最高的词，将其作为下一个生成的词。筛选的过程会一直持续，直到累计概率质量超过top_p或者概率质量最高的词已经全部选出为止。这样可以保证生成的文本符合上下文信息，同时又具有一定的随机性和多样性。

举个例子，假设当前上下文为“我想去看”，生成下一个词时，ChatGPT会计算每个可能的下一个词的条件概率分布，并根据top_p参数筛选出概率质量最高的若干个词，例如“电影”，“球赛”，“音乐会”等等。然后，ChatGPT会从筛选出的若干个词中随机选择一个词作为下一个生成的词，从而实现了多样性和可控性的平衡。

需要注意的是，top_p参数的取值通常在0.1到0.9之间，越大则生成文本的多样性越高，但也可能导致生成文本的不连贯和不合理。因此，在使用ChatGPT生成文本时，需要根据具体情况灵活调整top_p参数以达到最佳效果。

### n
n（也称为生成文本的长度），它可以用来控制生成文本的长度。

具体来说，n参数是一个整数，它指定了生成文本的长度，即生成的文本包含的单词数量。例如，如果n=10，则生成的文本将包含10个单词。需要注意的是，n参数的取值应该根据具体的应用场景进行调整，以便生成适当长度的文本。

在ChatGPT中，可以使用n参数来生成指定长度的文本，也可以使用n参数来生成不定长的文本。当n参数设置为一个较大的数值时（例如1000），ChatGPT会一直生成文本直到达到指定的长度或达到一定的停止条件（例如生成了结束标记或达到了最大生成步数）。这种方法可以生成长度不确定的文本，同时也可以避免生成文本长度过短或过长的问题。

需要注意的是，在使用n参数生成文本时，还需要结合其他的API参数（例如temperature和top_p）进行调整，以便生成符合要求的文本。例如，可以通过调整temperature参数来控制生成文本的多样性，通过调整top_p参数来控制生成文本的可控性，从而得到高质量的生成文本。

### max_tokens
用于控制生成的文本的长度。它表示生成文本时模型最多生成的token数量，其中一个token可以是一个词、一个标点符号或者一个空格等。例如，如果将max_tokens设置为100，那么生成的文本长度就不会超过100个token。

在实际应用中，我们可以根据需要设置max_tokens参数的值，以生成符合要求的文本长度。例如，在生成摘要或标题等短文本时，我们可以将max_tokens设置为一个较小的值，以避免生成过长的文本。而在生成长篇文章或者故事等较长文本时，则可以将max_tokens设置为较大的值，以获得更完整的文本内容。

需要注意的是，max_tokens的取值需要考虑模型的计算能力和内存限制。如果max_tokens设置过大，可能会导致模型计算时间和内存消耗过大，从而影响生成效率和稳定性。因此，在使用GPT模型时，我们需要根据具体应用场景和硬件条件，合理设置max_tokens参数的值。

### presence_penalty
presence_penalty的作用是增加模型在生成文本时的多样性。当presence_penalty为0时，模型会尽可能的重复或者延续给定的文本内容，生成的文本可能会比较单一或者缺乏创新性。而当presence_penalty不为0时，模型会尝试避免生成与给定文本重复或类似的内容，生成的文本将更加多样化和创新性。

需要注意的是，presence_penalty的取值需要根据具体应用场景进行调整。如果希望生成的文本与给定的文本内容尽可能相似，可以将presence_penalty设置为0；如果希望生成的文本内容与给定文本尽可能不同，可以将presence_penalty设置为较大的值。同时，presence_penalty的取值也需要与其他参数一起调整，以获得满足需求的生成文本结果。

### frequency_penalty
用于控制生成的文本中是否出现重复的文本片段。它的取值范围为[0, 2]，默认值为0。如果将frequency_penalty设置为非零值，那么生成的文本中会尽量避免出现重复的文本片段。

具体来说，frequency_penalty的作用是增加模型在生成文本时的多样性和独创性。当frequency_penalty为0时，模型可能会重复生成某些短语或者句子，导致生成的文本比较单一或者缺乏创新性。而当frequency_penalty不为0时，模型会尽量避免重复生成相同或类似的短语或者句子，生成的文本将更加多样化和创新性。

需要注意的是，frequency_penalty的取值也需要根据具体应用场景进行调整。如果希望生成的文本中不出现重复的短语或者句子，可以将frequency_penalty设置为较大的值；如果不考虑文本中的重复，可以将frequency_penalty设置为0。同时，frequency_penalty的取值也需要与其他参数一起调整，以获得满足需求的生成文本结果。




### user
代表您的最终用户的唯一标识符，可以帮助 OpenAI 监控和检测滥用行为。



