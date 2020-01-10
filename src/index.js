import React from 'react';
import ReactDOM from 'react-dom';
// import './index.css';


// const scyaleNames = {
//     c:'Celsius',
//     f:'Fahrenheit'
// };

// function toCelsius(fahrenheit) {
//     return (fahrenheit - 32) * 5 / 9;
//   }
  
// function toFahrenheit(celsius) {
//     return (celsius * 9 / 5) + 32;
//   }

//   function tryConvert(temperature,convert){
//     const input = parseFloat(temperature);
//     if(Number.isNaN(input)){
//         return '';
//     }
//     const output = convert(input);
//     const rounder = Math.round(output * 1000)/1000;
//     return rounder.toString();
// }

// function BoilingVerdict(props){
//     if(props.celsius>=100){
//         return <p>The water would boil.</p>
//     }
//     return <p>The water would not boil.</p>
// }

// class TemperatureInput extends React.Component{
//     constructor(props){
//         super(props);
//         this.handerChange=this.handerChange.bind(this);
//     }

//     handerChange(e){
//         this.props.onTemperatureChange(e.target.value);
//     }

//     render(){
//         const temperature = this.props.temperature;
//         const scale = this.props.scale;
//         return(
//             <fieldset>
//                 <legend>Enter temperature in {scaleNames[scale]}</legend>
//                 <input
//                     value={temperature}
//                     onChange={this.handerChange}
//                 />
//             </fieldset>
//         )
//     }
// }
// class Calculator extends React.Component{
//     constructor(props){
//         super(props);
//         this.state={
//             temperature: '',
//             scale: 'c'
//         }
//         this.handleFahrenheitChange=this.handleFahrenheitChange.bind(this)
//         this.handleCelsiusChange=this.handleCelsiusChange.bind(this)
//     }
//     handleCelsiusChange(temperature) {
//         this.setState({
//             scale: 'c',
//             temperature: temperature
//         })
//     }

//     handleFahrenheitChange(temperature){
//         this.setState({
//             scale:'f',
//             temperature
//         })
//     }
//     render(){
//         const scale = this.state.scale;
//         const temperature = this.state.temperature;
//         const celsius = scale==='f'?tryConvert(temperature,toCelsius) : 
//         temperature;
//         const fahrenheit = scale==='c'?tryConvert(temperature,toFahrenheit) : 
//         temperature;
//         return(
//             <div>
//                 <TemperatureInput scale="c"
//                  temperature={celsius}
//                  onTemperatureChange={this.handleCelsiusChange}
//                  />
//                 <TemperatureInput scale="f"
//                 temperature={fahrenheit}
//                 onTemperatureChange={this.handleFahrenheitChange}
//                 />
//                 <BoilingVerdict celsius={parseFloat(celsius)}/>
//             </div>
//         )
//     }
// }


// function FancyBorder(props){
//     return(
//         <div className={'FancyBorder FancyBorder'+props.color}>
//             {props.children}
//         </div>
//     )
// }
// function Dialog(props){
//     return(
//         <FancyBorder color="blue">
//             <h1 className="Dialog-title">
//                 {props.title}
//             </h1>
//             <h2 className="Dialog-message">
//                 {props.message}
//             </h2>
//             {props.children}
//         </FancyBorder>
//     )
// }

// function WelcomeDialog(){
//     return(
//       <Dialog
//         title="Welcome!"
//         message="This is a Message！！！"
//       />
//     )
// }
// class SignUpDialog extends React.Component{
//     constructor(props){
//         super(props);
//         this.state={login: ""}
//         this.handleChange=this.handleChange.bind(this);
//         this.handleSignUp=this.handleSignUp.bind(this);
//     }

//     render(){
//         return(
//             <Dialog title="This is a title"
//             message="This is a message">
//             <input 
//             value={this.state.login}
//             onChange={this.handleChange}
//             />
//             <button onClick={this.handleSignUp}>
//                     Click Me!
//             </button>
//             </Dialog>
//         )
//     }
//     handleChange(e){
//         this.setState({
//             login: e.target.value
//         })
//     }

//     handleSignUp(){
//         alert(`Welcome aboard, ${this.state.login}!`);
//     }
// }


// function Contacts() {
//     return <div className="Contacts"/>
// }
// function Chat() {
//     return <div className="Chat"/>
// }
// function SplitPane(props){
//     return(
//         <div className="SplitPane">
//             <div className="SplitPane-left">
//                 {props.left}
//             </div>
//             <div className="SplitPane-right">
//                 {props.right}
//             </div>
//         </div>
//     )
// }

// function App(){
//     return(
//         <SplitPane 
//         left={<Contacts/>} 
//         right={<Chat/>}
//         />
//     )
// }

//为每一个产品类别展示标题
// class ProductCategoryRow extends React.Component {
//     render() {
//       const category = this.props.category;
//       return (
//         <tr>
//           <th colSpan="2">
//             {category}
//           </th>
//         </tr>
//       );
//     }
//   }
//   //每一行展示一个产品
//   class ProductRow extends React.Component {
//     render() {
//       const product = this.props.product;
//       const name = product.stocked ?
//         product.name :
//         <span style={{color: 'red'}}>
//           {product.name}
//         </span>;
  
//       return (
//         <tr>
//           <td>{name}</td>
//           <td>{product.price}</td>
//         </tr>
//       );
//     }
//   }
//   //展示数据内容并根据用户输入筛选结果
//   class ProductTable extends React.Component {
//     render() {
//       const filterText = this.props.filterText;
//       const inStockOnly = this.props.inStockOnly;

//       const rows = [];
//       let lastCategory = null;
      
//       this.props.products.forEach((product) => {
//         if(product.name.indexOf(filterText)=== -1){
//             return;
//         }
//         if(inStockOnly && !product.stocked){
//             return;
//         }
//         if (product.category !== lastCategory) {
//           rows.push(
//             <ProductCategoryRow
//               category={product.category}
//               key={product.category} />
//           );
//         }
//         rows.push(
//           <ProductRow
//             product={product}
//             key={product.name} />
//         );
//         lastCategory = product.category;
//       });
  
//       return (
//         <table>
//           <thead>
//             <tr>
//               <th>Name</th>
//               <th>Price</th>
//             </tr>
//           </thead>
//           <tbody>{rows}</tbody>
//         </table>
//       );
//     }
//   }
//   //接受所有的用户输入
//   class SearchBar extends React.Component {
//       constructor(props){
//           super(props);
//         this.handleFilterTextChange=this.handleFilterTextChange.bind(this)
//         this.handleInStockChange=this.handleInStockChange.bind(this);
//       }
//       handleFilterTextChange(e){
//           this.props.onFilterTextChange(e.target.value);
//       }
//       handleInStockChange(e){
//           this.props.onInStockChange(e.target.value)
//       }
//     render() {
//       return (
//         <form>
//           <input type="text" 
//           placeholder="Search..."
//           value={this.props.filterText} 
//           onChange={this.handleFilterTextChange}
//           />
//           <p>
//             <input 
//             type="checkbox" 
//             checked={this.props.inStockOnly}
//             onChange={this.handleInStockChange}
//             />
//             {' '}
//             Only show products in stock
//           </p>
//         </form>
//       );
//     }
//   }
//   //是整个示例应用的整体
//   class FilterableProductTable extends React.Component {
//       constructor(props){
//           super(props);
//           this.state={
//               filterText:'',
//               inStockOnly: false
//           }
//           this.handleFilterTextChange=this.handleFilterTextChange.bind(this);
//           this.handleInStockChange=this.handleInStockChange.bind(this);
//       }
//       handleFilterTextChange(filterText){
//           this.setState({
//               filterText
//           })
//       }
//       handleInStockChange(inStockOnly){
//           this.setState({
//               inStockOnly
//           })
//       }
//     render() {
//       return (
//         <div>
//           <SearchBar 
//           filterText={this.state.filterText}
//           inStockOnly={this.state.inStockOnly}
//           onFilterTextChange={this.handleFilterTextChange}
//           onInStockChange={this.handleInStockChange}
//           />
//           <ProductTable 
//           products={this.props.products}
//           filterText={this.state.filterText}
//           inStockOnly={this.state.inStockOnly} />
//         </div>
//       );
//     }
//   }
  
  
//   const PRODUCTS = [
//     {category: 'Sporting Goods', price: '$49.99', stocked: true, name: 'Football'},
//     {category: 'Sporting Goods', price: '$9.99', stocked: true, name: 'Baseball'},
//     {category: 'Sporting Goods', price: '$29.99', stocked: false, name: 'Basketball'},
//     {category: 'Electronics', price: '$99.99', stocked: true, name: 'iPod Touch'},
//     {category: 'Electronics', price: '$399.99', stocked: false, name: 'iPhone 5'},
//     {category: 'Electronics', price: '$199.99', stocked: true, name: 'Nexus 7'}
//   ];
 const myh1 = React.createElement('h1',{id:'myh1',title:'this is a h1'},'这是一个H1')
// ReactDOM.render(<FilterableProductTable products={PRODUCTS}/>, document.getElementById('root'));
class MyH1 extends React.Component{
    render(){
        return(
            <div>
               myh1
            </div>
        )
    }
}
ReactDOM.render(<MyH1/>,document.getElementById('root'))

