---
title: "Query Operators"
weight: 20
date: 2023-07-29T12:22:42-05:00
---

EnSQL supports a subset of standard operators that are defined in SQL:1999. Operators include comparison operators and logical/boolean operators.

### Comparison Operators

Comparison operators are used to evaluate an expression that is composed of a left-side (e.g. `a` in the examples below), the operator, and the right-side (e.g. `b` in the examples below). These operators are typically used in the `WHERE` clause of a query for filtering results returned in an Ensign stream.

| Operator |        Syntax       | Description                             |
|:--------:|:-------------------:|-----------------------------------------|
|    `=`   |       `a = b`       | `a` is equal to `b`                     |
|   `!=`   |       `a != b`      | `a` is not equal to `b`                 |
|   `<>`   | `a <> b`            | `a` is not equal to `b` (alternate)     |
|    `>`   |       `a > b`       | `a` is greater than `b`                 |
|   `>=`   |       `a >= b`      | `a` is greater than or equal to `b`     |
|    `<`   |       `a < b`       | `a` is less than `b`                    |
|   `<=`   |       `a <= b`      | `a` is less than or equal to `b`        |
| `like`   | `a like 'pattern'`  | The pattern `'pattern'` is found in `a` |
| `ilike`  | `a ilike 'pattern'` | Case-insensitive `like` search          |

### Logical/Boolean Operators

Logical operators return the result of a Boolean operation on an input expression that is composed of a left-side (e.g. `a` in the examples below), the operator, and the right-side (e.g. `b` in the examples below). Both input expressions on the left and ride side must evaluate to a boolean value.

Logical operators can only be used as a predicate/condition e.g. in the `WHERE` clause of a SQL statement.

| Operator |   Syntax  | Description                                                              |
|:--------:|:---------:|--------------------------------------------------------------------------|
|   `AND`  | `a AND b` | Both `a` and `b` must evaluate to `true` to be `true` otherwise `false`  |
|   `OR`   |  `a OR b` | Either `a` or `b` must evaluate to `true` to be `true` otherwise `false` |

The order of precedence of these operators is shown below from highest to lowest:

1. AND
2. OR

**NOTE**: When we add the `NOT` logical operator in the future, it will have the highest precedence.