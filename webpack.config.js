const path = require('path')
const HtmlWebPackPlugin =require('html-webpack-plugin') //导入 在内存中自动生成 index页面的插件


//创建一个插件的实例对象
const htmlPlugin = new HtmlWebPackPlugin({
    template:path.join(__dirname,'/public/index.html'), //用于生成的源文件 path.join地址拼接  __dirname表示当前文件的地址 
    filename: 'index.html'//生成的内存中的首页的名称
})


module.exports = {
    mode: 'development',
    plugins:[
        htmlPlugin
    ],
    module:{
        rules: [
           { test: /\.js|jsx$/,use :'babel-loader',exclude: /node_modules/ } //exclude是排除项
        ]
    },
    resolve:{
        extensions:['.js','.jsx','.json'],//表示这几个文件的后缀名可以不写
        alias: {
            '@':path.join(__dirname,'./src')//之后使用@就表示项目根目录下的src
        }
    }
}