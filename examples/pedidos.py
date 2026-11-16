# Pedidos de una librería: clases con herencia, listas, tuplas y diccionarios.


class Producto:
    def __init__(self, nombre, precio):
        self.nombre = nombre
        self.precio = precio

    def precio_final(self):
        return self.precio


class ProductoEnOferta(Producto):
    def __init__(self, nombre, precio, descuento):
        Producto.__init__(self, nombre, precio)
        self.descuento = descuento

    def precio_final(self):
        return self.precio * (100 - self.descuento) // 100


class Pedido:
    cantidad_de_pedidos = 0

    def __init__(self, cliente):
        self.cliente = cliente
        self.lineas = []
        Pedido.cantidad_de_pedidos = Pedido.cantidad_de_pedidos + 1

    def agregar(self, producto, cantidad):
        self.lineas.append((producto, cantidad))

    def total(self):
        suma = 0
        for linea in self.lineas:
            producto, cantidad = linea
            suma = suma + producto.precio_final() * cantidad
        return suma


catalogo = {}
catalogo["quijote"] = Producto("Don Quijote", 12000)
catalogo["ficciones"] = ProductoEnOferta("Ficciones", 9000, 20)
catalogo["rayuela"] = Producto("Rayuela", 10500)
stock = {"quijote": 3, "ficciones": 1, "rayuela": 1}

encargos_de_ana = [("quijote", 2), ("ficciones", 1)]
encargos_de_beto = [("ficciones", 1), ("rayuela", 1), ("borges", 1)]

pedidos = []
for cliente_y_encargos in [("Ana", encargos_de_ana), ("Beto", encargos_de_beto)]:
    cliente, encargos = cliente_y_encargos
    pedido = Pedido(cliente)
    for encargo in encargos:
        clave, cantidad = encargo
        if clave not in catalogo:
            print(cliente, "pidió algo que no vendemos:", clave)
        elif stock.get(clave, 0) < cantidad:
            print(cliente, "se quedó sin", catalogo[clave].nombre)
        else:
            stock[clave] = stock[clave] - cantidad
            pedido.agregar(catalogo[clave], cantidad)
    pedidos.append(pedido)

for pedido in pedidos:
    nombres = []
    for linea in pedido.lineas:
        nombres.append(linea[0].nombre)
    print(pedido.cliente, nombres, "total:", pedido.total())

agotados = []
for clave in stock:
    if stock[clave] == 0:
        agotados.append(clave)
print("agotados:", agotados)
print("stock:", stock)

historial = pedidos
historial.append(Pedido("Caro"))
print("pedidos creados:", Pedido.cantidad_de_pedidos, "y la lista original tiene", len(pedidos))

ultimos = pedidos[-2:]
primero, segundo = ultimos
print("últimos clientes:", (primero.cliente, segundo.cliente))
clientes = []
for pedido in pedidos:
    clientes.append(pedido.cliente)
print("al revés:", clientes[::-1])
