# dailyTerminal

dailyTerminal es un pequeño script hecho en Go que imprime qué se celebra el día de hoy, pero solo en la primera terminal que abras cada día. Las siguientes terminales del mismo día no vuelven a imprimirlo — pensado para integrarse a tu .bashrc junto con lo que ya tengas ahí.

## Cómo funciona
 
1. Al ejecutarse, revisa un archivo de caché (`~/.cache/dayfact/last_shown`) donde guarda la última fecha en que ya mostró algo.
2. Si esa fecha es **hoy**, no imprime nada y termina con código de salida `1`.
3. Si no, consulta la API pública de Wikipedia ([`onthisday/holidays`](https://en.wikipedia.org/api/rest_v1/feed/onthisday/holidays/07/25)), imprime las festividades del día, guarda la fecha de hoy en el caché, y termina con código `0`.
Ese código de salida (`0` o `1`) es justo lo que usa `.bashrc` para decidir si muestra la efeméride o cae a tu quote de siempre.
 
## Requisitos
 
- Go 1.22 o superior (solo para compilar; el binario ya compilado no necesita Go instalado).
- Conexión a internet al momento de correrlo (si no hay, falla en silencio y no rompe la terminal).

---
## Instalación

Clonar el repo:
```bash
git clone https://github.com/Penya123/dailyTerminal.git
cd dailyTerminal
```

compila e instala el binario en `~/go/bin`:
 
```bash
go install .
```
 
Asegúrate de que `~/go/bin` esté en tu `PATH`. Si no lo está, agrega esto a tu `.bashrc`:
 
```bash
export PATH="$HOME/go/bin:$PATH"
```
 
> **¿Vas a usarlo en otra máquina o arquitectura distinta** (ej. un celular ARM, Raspberry Pi)? No hace falta tener Go instalado ahí: compila desde tu máquina con variables `GOOS`/`GOARCH` y copia el binario resultante. Ejemplo para ARM64:
> ```bash
> GOOS=linux GOARCH=arm64 go build -o dailyTerminal-arm64 .
> ```
 
## Uso
 
Agrega esto a tu `.bashrc`, después de lo que ya imprimes al abrir terminal (ej. `neofetch`):
 
```bash
if dailyTerminal; then
    : # ya se imprimieron las celebraciones del día, no hace falta nada más
else
    echo "Nuena Terminal" | lolcat   # tu comportamiento normal
fi
```
 
Abre una terminal nueva y deberías ver algo como:
 
```
Hoy es 25 de julio, se celebra:
  • Guanacaste Day (Costa Rica)
  • National Day of Galicia (Galicia, Spain)
  • Tenjin Matsuri (Osaka, Japan)
```
 
Abre una segunda terminal el mismo día: no se imprime nada de esto, y en su lugar corre el `else` (tu frase de siempre).
 
## Notas
 
- Los nombres de las festividades vienen tal cual los devuelve Wikipedia en inglés — no hay traducción automática al español disponible en esta API.
- Se omiten intencionalmente las entradas de "Christian feast day" (festividades religiosas cristianas), que suelen ser la mayoría de los resultados y no aportan mucho al propósito del script.
- Si por alguna razón el caché no se puede leer/escribir, o no hay internet, el programa simplemente no imprime nada y sale con código `1` — nunca debería romper o colgar tu terminal.
