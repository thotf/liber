# liber
使用golang开发的解释器

![image](https://user-images.githubusercontent.com/48403565/131203212-84726d8a-02cc-4422-8675-798fc5fa39c2.png)

使用一遍式扫描。
没有为每一阶段单独生成语法树，每次进行语法分析时调用词法分析，同时进行语义分析。

定义有 string、int、float、true、false、fun、list  几种数据类型，多种类型之间可相互转换
list的简易的lua 表的实现。

定义函数时要注意函数中每个return后都要用同样数量的返回值。

abc.txt 为冒泡排序测试代码

支持调用golang函数，需要先定义gofuner接口。目前已经实现的有fmt.Println()、len() 的实现

尚未实现一些符合逻辑的简单语法比如++、+= ...
尚未支持UPVALUE 只是留下了接口。
循环控制语句目前只有while，忘记添加break和continue了。所幸这只需要一点简单的写回技术就可实现
控制语句if、else 已实现。

语法分析采用TDOP算法，这更像是算符优先算法的改良版。和算符优先比起来，不需要通过定义文法来计算各种符号的优先级，而是直接定义好优先级，再进行比较的算法
