# Color Counter

[Artículo en Medium](https://medium.com/@umarquez_mx/10-100-100daysofc0d3-7187417844ad)

Obtiene la frecuencia de colores de una imagen iterando sobre cada pixel que la compone y escribiendo los datos del proceso en un archivo json, esto se logra utilizando una tabla hash para clasificar cada color y un árbol binario de búsqueda para realizar el ordenamiento por frecuencia.

## Uso

### `-img <imagen/de/entrada.jpg>`
Establece la ruta de la imagen a procesar, si se omite el sistema descarga una imagen aleatoria.

Ejemplo: `-img ./example.jpg`

### `-out <archivo/de/salida.json>`
Establece la ruta del archivo de salida, si se omite el sistema utilizará el mismo nombre de la imagen de entrada y concatenará la extensión `.json`

Ejemplo: `-out ./example.json`

`#100DayOfC0D3`