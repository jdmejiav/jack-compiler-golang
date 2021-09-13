# jack-compiler-golang

Compilador de jack a VM Programado en Golang, sin usar code generators.

Jack language compiler from jack to VM. Without code generators


# status del proyecto // Project Status:
 - **Tokenyze:** 100%:
 
 The analyzer don't recognize coment in the  \/\* \*\/ way. example:
 
 /\*
 
 \<Comment\>

 */ 
 
 **Only in the way:**
 
 //\<comment\>
 
 - **Parser:** 100% Only accepts one variable definition per code line. Example:
 
 **Wrong:**

    var int a,b,c,d;

**Correct:**

    var int a;
    var int b;
    var int c;
    var int d;


 - **Code generator:** 50% (Currently i only compiles fucntion definition, class attributes definition and Do statements)


# How to Use:


It is necessary to intall golang:

## On Windows / MAC:
You can go to the official website and download the installer
<a href="https://golang.org/">Go</a>

## On Linux:
Usually it comes installed on all the distributions, but if it isn't you have to look how to install the package on you specific distribution, here some examples

### Ubuntu 
    $ apt-get install golang-go
### Arch
    $ pacman -S go


for the program execution, it is neccesary to compile first, for this, go to the src folder inside the project and compile the code in the following way:

    $ cd src
    $ go build main.go token.go  pila.go  parser.go  linkedList.go
 
Once is compilled, run the ./main executable file and pass as parameter the route of the .jack file you want to compile (on the /JackTest/ directory are some .jack code examples)

    $ ./main.exe ../jackText/Main.jack  
    $ ./main.exe ../jackTest/Barra.jack
    $ ./main.exe ../jackText/Bola.jack

