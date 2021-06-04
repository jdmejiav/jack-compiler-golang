# jack-compiler-golang

Compilador de jack a VM Programado en Golang

# status del proyecto:
 - **Tokenizado:** 100% No acepta comentarios de la forma:
 
 /\*\
 
 <Comentario\>\

 */ 
 
 **únicamente de la forma:**
 
 //\<comentario\>
 
 - **Parser:** 100% No acepta definiciones de más de una variable en una sola línea ej:
 
 **Incorrecto: **

    var int a,b,c,d;

**Correcto: **

    var int a;

    var int b;

    var int c;

    var int d;


 - **Code generator:*** 50% (actualmente solo compila, defición de funciones, definición de atributos de clase y do statements)


# Modo de Uso:

Es necesario instalar golang, para hacerlo se puede hacer desde la página oficial
[Instalador Golang](https://golang.org/)

