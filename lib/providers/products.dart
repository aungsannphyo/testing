import 'dart:convert';

import 'package:flutter/material.dart';
import 'product.dart';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import '../models/http_exception.dart';
import 'package:dio/dio.dart';

class Products with ChangeNotifier {
  List<Product> _items = [];
  final String _authToken;
  final String _userId;

  Products(this._authToken, this._userId, this._items);

  final String _url = 'testing-shop-75ec2-default-rtdb.firebaseio.com';

  List<Product> get items {
    return [..._items];
  }

  List<Product> get favoriteItems {
    return _items.where((prodItem) => prodItem.isFavorite).toList();
  }

  Future<void> fetchAndSetProducts() async {
    // var url = Uri.https(_url, '/products.json', {
    //   'auth': _authToken,
    //   'orderBy': "creatorId",
    //   'equalTo': _userId,
    // });

    final Uri uri = Uri.parse(
        'https://testing-shop-75ec2-default-rtdb.firebaseio.com/products.json?auth=$_authToken&orderBy="creatorId"&equalTo="$_userId"');

    var favoriteUrl = Uri.https(_url, '/userFavorites/$_userId.json', {
      'auth': _authToken,
    });

    try {
      final response = await Dio().getUri(uri);
      final Map<String, dynamic>? data = response.data;

      final List<Product> loadedProducts = [];
      if (data == null) {
        return;
      }
      final favoriteResponse = await http.get(favoriteUrl);
      final favoriteData = json.decode(favoriteResponse.body);
      data.forEach((productId, productData) {
        loadedProducts.add(
          Product(
            id: productId,
            title: productData['title'],
            description: productData['description'],
            price: productData['price'],
            imageUrl: productData['imageUrl'],
            isFavorite:
                favoriteData == null ? false : favoriteData[productId] ?? false,
          ),
        );
      });
      _items = loadedProducts;
      notifyListeners();
    } catch (error) {
      rethrow;
    }
  }

  Future<void> addProduct(Product product) async {
    var url = Uri.https(_url, '/products.json', {'auth': _authToken});

    try {
      await http.post(url,
          body: json.encode({
            'title': product.title,
            'description': product.description,
            'price': product.price,
            'imageUrl': product.imageUrl,
            'creatorId': _userId
          }));
      final newProduct = Product(
        title: product.title,
        description: product.description,
        price: product.price,
        imageUrl: product.imageUrl,
        id: DateTime.now().toString(),
      );

      _items.add(newProduct);
      notifyListeners();
    } catch (error) {
      rethrow;
    }
  }

  Product findById(String id) {
    return _items.firstWhere((prod) => prod.id == id);
  }

  Future<void> updateProduct(String id, Product newProduct) async {
    var url = Uri.https(_url, '/products/$id.json', {'auth': _authToken});

    final productIndex = _items.indexWhere((prod) => prod.id == id);
    if (productIndex >= 0) {
      try {
        await http.patch(url,
            body: json.encode({
              'title': newProduct.title,
              'description': newProduct.description,
              'price': newProduct.price,
              'imageUrl': newProduct.imageUrl,
            }));
        _items[productIndex] = newProduct;
      } catch (error) {
        rethrow;
      }
    } else {
      return;
    }

    notifyListeners();
  }

  Future<void> deleteProduct(String id) async {
    var url = Uri.https(_url, '/products/$id.json', {'auth': _authToken});
    final existingProductIndex = _items.indexWhere((prod) => prod.id == id);
    Product? existingProduct = _items[existingProductIndex];
    _items.removeAt(existingProductIndex);
    notifyListeners();

    final respone = await http.delete(url);

    if (respone.statusCode >= 400) {
      _items.insert(existingProductIndex, existingProduct);
      notifyListeners();
      throw HttpException('Could not delete product.');
    }
    existingProduct = null;
  }
}
