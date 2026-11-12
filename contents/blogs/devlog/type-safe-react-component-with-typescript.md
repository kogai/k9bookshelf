---
title: "TypeScript✕ReactJSで型安全なComponentを手に入れる"
published_at: "2016-08-10T23:23:35+09:00"
tags:
  - "javascript"
  - "typescript"
  - "reactjs"
---

## React.Proptypesに感じている不満

React組み込みのProptypesによる型検査は、実行時検証であり、アプリケーションを起動して、実際にコンポーネントを描画してみるまでエラーに気づけない。
また、React独自の機能であるので、他の下層DOMライブラリでも同様の型検査をできるわけではない。
さらに、Proptypesはコンポーネントに定義するので、型定義がViewに紐付いてしまう(ドメイン知識がViewに紐付いてしまう)
型定義の多くはドメイン知識に紐付いており、ViewではなくModel(State)レイヤーで定義したい。

そこで本稿ではTypeScriptの静的解析の支援を受けて、型安全なコンポーネントを定義することでこの問題を解決する方法について説明したい。
※なお、ReactにはFacebook謹製のFlowTypeがあるが、あくまでアプリケーションの主体はModel層にあり、Model層はRxJSによるState川を作る設計にしたいという理由でTypeScriptを選択している

## 導入方法

まず、TypeScriptとReactをインストールする。
※TypeScriptは2.x系のバージョンがベータ版としてリリースされているので、そちらを使うのであればtypingsは不要。

```bash
npm i -S typescript react typings
```

次にReactの型定義ファイルを取得します

```bash
$(npm bin)/typings install -S react
# TypeScript@2の場合は npm install -D @types/react
```

これで準備は完了。

## コンポーネントの定義

それでは実際にコンポーネントを定義する。
ReactにはReactComponentと、StatelessComponentという2種のコンポーネント型がある。

ReactComponentは、一般的なReactのコンポーネントで、React.Componentを継承して定義するコンポーネント。
StatelessComponentは、関数として定義できるコンポーネントで、(当然だが)ライフサイクルメソッド、Stateを持たない。
Viewにロジックが入り込む余地をなくすことができるので、私は特に理由がない限りこちらのコンポーネントを使うようにしている。

さて、ReactComponentとStatelessComponentのコンポーネント定義例を作成する前に、それぞれの型定義を確認したい。

```typescript
// ReactComponent
class Component<P, S> implements ComponentLifecycle<P, S> {
  constructor(props?: P, context?: any);
  setState(f: (prevState: S, props: P) => S, callback?: () => any): void;
  setState(state: S, callback?: () => any): void;
  forceUpdate(callBack?: () => any): void;
  render(): ReactElement<any>;
  
  props: P & { children?: ReactNode };
  state: S;
  context: {};
  refs: {
    [key: string]: ReactInstance
  };
}
  
// StatelessComponent
interface StatelessComponent<P> {
  (props?: P, context?: any): ReactElement<any>;
  propTypes?: ValidationMap<P>;
  contextTypes?: ValidationMap<any>;
  defaultProps?: P;
  displayName?: string;
}
```

ReactComponentは`P`,`S`、StatelessComponentは`P`という型変数を受け取ることに注目したい。
これはPropsとStateを表す型変数で、それぞれコンポーネントのPropsとStateに受け取れる型を定義できる。
※余談だが、コンポーネント階層に暗黙に流れていくcontextというオブジェクトもあるが、こちらはany型になっている。

つまり、PropTypesで実行時に検証していたPropsの型を、TypeScriptのジェネリクス機能によって、静的に解析・検証ができるというだ。
これによって、Model(Store)レイヤーで定義したインターフェースに変更があった場合でも、アプリケーションを実行することなく変更を検出・コンポーネントの定義を変更でき、非常に快適なコンポーネント作成を行うことができる。

## 作例

では実際にコンポーネントの作例を書いてみたい。
例によってTodoアプリケーションを想定して、ごく簡単なTodoリストを描画したいとする。
(なお、JSX記法を用いて仮想DOMを定義しているので、拡張子を*.tsxとしている)

```typescript
import { StatelessComponent, Component } from "react";
  
// Propsの型を定義する
interface ITodo {
  label: string;
  isCompleted: boolean;
}
  
interface ITodos {
  todos: ITodo[];
}
  
// StatelessComponentの型引数に上で定義した型を渡す
const Todo: StatelessComponent<ITodos> = ({ label, isCompleted }) => (
  <li className={ isCompleted ? "Todo--is-complete" : "Todo" }>{ label }</li>
);
  
const Todos: StatelessComponent<ITodoProps> = ({ todos }) => (
  <ul>
    {todos.map(todo => <Todo { ...todo }/>)}
  </ul>
);
  
// ReactComponentとして定義するなら...
class Todo extends Component<ITodos, void> {
  render() {
    const { label, isCompleted } = this.props;
    return <li className={ isCompleted ? "Todo--is-complete" : "Todo" }>{ label }</li>
  }
}
  
class Todos extends Component<ITodos, void> {
  render() {
    const { todos } = this.props;
    return (
      <ul>
        {todos.map(todo => <Todo { ...todo }/>)}
      </ul>
    )
  }
}
```

この作例では必要なかったが、React.Componentの第二型引数に渡した型で、this.stateの型も定義できる。
※個人的にはそもそもコンポーネント内部でStateを操作しないケースの方が多いと感じているが、アプリケーションのサイズに依存するだろう

以上で型安全な子コンポーネントが作成できた。
これらのコンポーネントの親となるルートコンポーネントでルートStateを購読し、下層のコンポーネントに流してあげることで、アプリケーションのViewとしての定義が完成する。

```typescript
export class RootComponent extends Component<void, ITodos> {
  componentWillMount() {
    rootModel.subscribe(state => this.setState(state));
  }
  
  render() {
    return <Todos { ...this.state } />;
  }
}
```

## まとめ

以上、TypeScriptでReactコンポーネントを作成することで型安全なコンポーネントを定義する方法を紹介してみた。
参考になればうれしい。
