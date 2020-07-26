import * as React from 'react'

const logo = require('../media/SDV_logo_3x.png')

export default function (props: { message: React.FunctionComponentElement<any> }) {
  return (
    <div className="App">
      <img src={logo} className="App-logo" alt="logo"/>
      {props.message}
    </div>
  )
}
