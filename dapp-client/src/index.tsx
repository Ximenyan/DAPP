import * as Ontology from 'ontology-dapi';
import * as React from 'react';
import * as ReactDOM from 'react-dom';

import * as base58 from 'bs58';
// import { BrowserRouter, Route } from 'react-router-dom';
// import { Home } from './home';
// import { Message } from './message';
// import { Network } from './network';
// import { Oep4 } from './oep4';
// import { Provider } from './provider';
// import { SmartContract } from './smartContract';

Ontology.client.registerClient({});
/*const App: React.SFC<{}> = () => (
  <BrowserRouter>
    <>
      <Route path="/" exact={true} component={Home} />
      <Route path="/network" exact={true} component={Network} />
      <Route path="/oep4" exact={true} component={Oep4} />
      <Route path="/smart-contract" exact={true} component={SmartContract} />
      <Route path="/message" exact={true} component={Message} />
      <Route path="/provider" exact={true} component={Provider} />
    </>
  </BrowserRouter>
);*/

function ab2hexstring(arr: any): string {
  let result: string = '';
  const uint8Arr: Uint8Array = new Uint8Array(arr);
  for (let i = 0; i < uint8Arr.byteLength; i++) {
    let str = uint8Arr[i].toString(16);
    str = str.length === 0
        ? '00'
        : str.length === 1
            ? '0' + str
            : str;
    result += str;
  }
  return result;
}
function base58ToHex(base58Encoded: string) {
  const decoded = base58.decode(base58Encoded);
  const hexEncoded = ab2hexstring(decoded).substr(2, 40);
  alert(decoded);
  // if (base58Encoded !== hexToBase58(hexEncoded)) {
  //   throw new Error('[addressToU160] decode encoded verify failed');
  // }
  return hexEncoded;
}
const Login: React.SFC<{}> = () => {
  async function onGetAccount() {
    const account = await Ontology.client.api.asset.getAccount();
    // alert('onGetPublicKey: ' + base58ToHex(account));
    ReactDOM.render(<>{account}</>, document.getElementById('address') as HTMLElement);
    alert('onGetAccount: ' + JSON.stringify(account));
  }
  return(
<>
    <a href="#" className={'ui-btn ui-corner-all ui-shadow ui-icon-home ui-btn-icon-left'}>注销</a>
    <h1 id="address">Ugly交易所</h1>
    <a href="#" className={'ui-btn ui-corner-all ui-shadow ui-icon-search ui-btn-icon-left'} onClick={onGetAccount}>
      登录
    </a>
</>
  );
};

function ActionLink() {
  function convertValue(value: string, type: Ontology.ParameterType) {
    switch (type) {
      case 'Boolean':
        return Boolean(value);
      case 'Integer':
        return Number(value);
      case 'ByteArray':
        return value;
      case 'String':
        return value;
    }
  }
  async function handleClick() {
    const account = await Ontology.client.api.asset.getAccount();
    const scriptHash: string = 'faf0bb2cb8b525477e8bbfdb4d81dd373b47c390';
    const operation: string = 'CreateBuyOrder';
    const gasPrice: number = 500;
    const gasLimit: number = 20000000;
    const requireIdentity: boolean = true;
    const hexstr: string = base58ToHex(account);
    const parametersRaw: any[] = [
    {type: 'ByteArray', value: hexstr },
    { type: 'String', value: '_ONG_ONT_' },
    { type: 'Integer', value: Number(1) },
    { type: 'Integer', value: Number(1000000) },
    { type: 'Integer', value: Number(0) },
    { type: 'Integer', value: Number(2) }];
    const args = parametersRaw.map((raw) => ({ type: raw.type, value: convertValue(raw.value, raw.type) }));
    try {
      const result = await Ontology.client.api.smartContract.invoke({
        scriptHash,
        operation,
        args,
        gasPrice,
        gasLimit,
        requireIdentity
      });
      // tslint:disable-next-line:no-console
      console.log('onScCall finished, result:' + JSON.stringify(result));
    } catch (e) {
      alert('onScCall canceled');
      // tslint:disable-next-line:no-console
      console.log('onScCall error:', e);
    }
  }
  return (
  <a href="#" onClick={handleClick} className={'ui-btn ui-btn-inline ui-corner-all'} id="sell_ont_btn">
    &emsp;挂单&emsp;
  </a>);
}
ReactDOM.render(<ActionLink />, document.getElementById('buy_ont_btn'));
ReactDOM.render(<Login />, document.getElementById('head') as HTMLElement);
// ReactDOM.render(<App />, document.getElementById('root') as HTMLElement);
