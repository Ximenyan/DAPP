import * as Ontology from 'ontology-dapi';
import * as React from 'react';
import * as ReactDOM from 'react-dom';

import * as base58 from 'bs58';

Ontology.client.registerClient({});

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
  return hexEncoded;
}

const Login: React.SFC<{}> = () => {
  async function onGetAccount() {
    const account = await Ontology.client.api.asset.getAccount();
    ReactDOM.render(<>{account}</>, document.getElementById('address') as HTMLElement);
    alert('onGetAccount: ' + JSON.stringify(account));
  }

  return (
    <>
      <a href="#" className={'ui-btn ui-corner-all ui-shadow ui-icon-home ui-btn-icon-left'}>注销</a>
      <h1 id="address">Ugly交易所</h1>
      <a href="#" className={'ui-btn ui-corner-all ui-shadow ui-icon-search ui-btn-icon-left'}
         onClick={onGetAccount}>
        登录
      </a>
    </>
  );
};

function CreateOrder(props: any) {
  const isSell: boolean = props.isSell;

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
    const scriptHash: string = 'ac1bef614a9cbd2ed44d86fd2a8341c2a3c7300d';
    let operation: string = 'CreateBuyOrder';
    const gasPrice: number = 500;
    const gasLimit: number = 20000000;
    const requireIdentity: boolean = true;
    const hexstr: string = base58ToHex(account);
    $.ajaxSettings.async = false;
    let preId: number = Number(0);
    let nextId: number = Number(0);
    let price: number = Number(0); // $('#buy_price_input').val()
    let amount: number = Number(0); // $('#buy_amount_input').val()
    let orderType: string = '_BUY___List_Tail_Order___ONG_ONT_';
    if (isSell) {
      orderType = '_SELL___List_Tail_Order___ONG_ONT_';
      operation = 'CreateSellOrder';
      price = Number($('#sell_price_input').val()) * 1000000000;
      amount = Number($('#sell_amount_input').val());
      // alert(price);
    } else {
      price = Number($('#buy_price_input').val()) * 1000000000;
      amount = Number($('#buy_amount_input').val());
      // alert(price);
    }
    // alert(operation);
    $.get('/api?req_type=create_order&order_type=' + orderType + '&price=' + price + '',
      function(data, status) {
        const numArr = JSON.parse(data);
        preId = numArr[0];
        nextId = numArr[1];
        // alert(preId);
        // alert(nextId);
      });
    const parametersRaw: any[] = [
      {type: 'ByteArray', value: hexstr},
      {type: 'String', value: '_ONG_ONT_'},
      {type: 'Integer', value: amount},
      {type: 'Integer', value: price},
      {type: 'Integer', value: Number(preId)},
      {type: 'Integer', value: Number(nextId)}];
    const args = parametersRaw.map((raw) => ({type: raw.type, value: convertValue(raw.value, raw.type)}));
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

  if (isSell) {
    return (
      <a href="#" onClick={handleClick} className={'ui-btn ui-btn-inline ui-corner-all'} id="sell_ont_btn">
        &emsp;卖出&emsp;
      </a>);
  }
  return (
    <a href="#" onClick={handleClick} className={'ui-btn ui-btn-inline ui-corner-all'} id="buy_ont_btn">
      &emsp;挂单&emsp;
    </a>);
}

ReactDOM.render(<CreateOrder isSell={false}/>, document.getElementById('buy_btn'));
ReactDOM.render(<CreateOrder isSell={true}/>, document.getElementById('sell_btn'));
ReactDOM.render(<Login/>, document.getElementById('head') as HTMLElement);
