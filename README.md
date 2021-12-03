## Excel-Tools

> 该工具的主要作用是用于Excel导出JSON、CSV、SQL格式，目前正在开发中。
> 

### 支持以下数据类型
- number 数字类型。
- bool 布尔。
- string 字符串。
- date 日期类型。
- object 对象，同JSON对象一致。
- array 数组，同JSON数组一致。
- id 主键类型(当表中有id类型时，json会以hash格式输出，否则以array格式输出)。

### 表头规则
- 基本数据类型(string,number,bool)时候，一般不需要设置会自动判断，但是也可以明确声明数据类型。
- 字符串类型：命名形式 列名#string 。
- 数字类型：命名形式 列名#number 。
- 日期类型：列名#date 。日期格式要符合标准日期格式。比如YYYY/M/D H:m:s or YYYY/M/D 等等。
- 布尔类型：命名形式 列名#bool 。
- 数组：命名形式 列名#[]。
- 对象：命名形式 列名#{}。
- 主键：命名形式列名#id ,表中只能有一列。
- 列名字以#开头则不导出此列。

### sheet规则
sheet名字以#开头则不导出此表。

