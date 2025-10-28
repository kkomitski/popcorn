# Parser Expression Precedence

This document describes the **expression precedence** implemented in the `Parser` class (`parser.ts`). Operator precedence determines the order in which expressions are parsed and evaluated, ensuring that the resulting Abstract Syntax Tree (AST) reflects the correct grouping and evaluation order.

## Precedence Order (Lowest to Highest)

The parser handles expressions by starting with the lowest precedence and recursively parsing higher-precedence expressions. Here is the order, from lowest to highest:

1. **Assignment Expression**  
   _Example:_ `let x = 5;`  
   _Parser method:_ `parseAssignmentExpression`

2. **Object Expression**  
   _Example:_ `{ key: value }`  
   _Parser method:_ `parseObjectExpression`

3. **Additive Expression**  
   _Example:_ `a + b - c`  
   _Parser method:_ `parseAdditiveExpression`

4. **Multiplicative Expression**  
   _Example:_ `a * b / c % d`  
   _Parser method:_ `parseMultiplicativeExpression`

5. **Call/Member Expression**  
   _Examples:_  
   - Function call: `add(2, 3)`  
   - Member access: `obj.prop` or `arr[0]`  
   _Parser methods:_ `parseCallMemberExpression`, `parseMemberExpression`, `parseCallExpression`

6. **Primary Expression**  
   _Examples:_  
   - Literals: `42`, `"hello"`  
   - Identifiers: `x`  
   - Grouping: `(a + b)`  
   _Parser method:_ `parsePrimaryExpression`

## How Precedence is Implemented

- The parser starts with `parseExpression()`, which calls `parseAssignmentExpression()`.
- Each parsing method delegates to the next higher-precedence method.
- For example, `parseAssignmentExpression` calls `parseObjectExpression`, which calls `parseAdditiveExpression`, and so on.
- This recursive descent structure ensures that higher-precedence operations are parsed first and grouped correctly in the AST.

## Example

For the expression:

```js
let x = a + b * c;