# jack-compiler-golang

Compilador de jack a VM Programado en Golang

# status del proyecto:
 - **Tokenizado:** 100% No acepta comentarios de la forma:
 
 /\*
 
 \<Comentario\>

 */ 
 
 **únicamente de la forma:**
 
 //\<comentario\>
 
 - **Parser:** 100% No acepta definiciones de más de una variable en una sola línea ej:
 
 **Incorrecto:**

    var int a,b,c,d;

**Correcto:**

    var int a;
    var int b;
    var int c;
    var int d;


 - **Code generator:*** 50% (actualmente solo compila, defición de funciones, definición de atributos de clase y do statements)


# Modo de Uso:

Es necesario instalar golang, para hacerlo se puede hacer desde la página oficial
[Instalador Golang](https://golang.org/)

Para ejecutar el programa, es necesario compilarlo, para esto, dirígase a la carpeta src dentro del proyecto y compile el código de la siguiete manera:
    
    $ cd src
    $ go build main.go token.go  pila.go  parser.go  linkedList.go

Una vez compilado, pase como parámetro un que desee compilar (En el directorio jackTest/ hay varios ejemplos de archivos .jack). Ej:

    $ ./main.exe ../jackText/Main.jack  
    $ ./main.exe ../jackTest/Barra.jack
    $ ./main.exe ../jackText/Bola.jack
