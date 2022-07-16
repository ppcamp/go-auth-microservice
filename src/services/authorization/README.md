# authorization

This part of the code will be responsible to check if the passed groups in the JWT token can access
an endpoint or not (RBAC). To achieve this we can use [CASBIN](https://casbin.org/) or
[indigo](github.com/ezachrisen/indigo)