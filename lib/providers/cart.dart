import 'package:flutter/widgets.dart';

import 'package:flutter/foundation.dart';

class CartItem {
  final String id;
  final String title;
  final int quantiry;
  final double price;

  CartItem(
      {required this.id,
      required this.title,
      required this.quantiry,
      required this.price});
}

class Cart with ChangeNotifier {
  Map<String, CartItem>? _items = {};

  Map<String, CartItem> get items {
    return {..._items!};
  }

  int get itemCount {
    return _items == null ? 0 : _items!.length;
  }

  double get totalAmount {
    var total = 0.0;
    _items!.forEach((key, cardItem) {
      total += cardItem.price * cardItem.quantiry;
    });
    return total;
  }

  void addItem(String productId, double price, String title) {
    if (_items!.containsKey(productId)) {
      //change quantity
      _items!.update(
        productId,
        (value) => CartItem(
            id: value.id,
            title: value.title,
            quantiry: value.quantiry + 1,
            price: value.price),
      );
    } else {
      _items!.putIfAbsent(
        productId,
        () => CartItem(
            id: DateTime.now().toString(),
            title: title,
            quantiry: 1,
            price: price),
      );
    }
    notifyListeners();
  }

  void removeItem(String productId) {
    _items!.remove(productId);
    notifyListeners();
  }

  void removeSingleItem(String productId) {
    if (!_items!.containsKey(productId)) {
      return;
    }
    if (_items![productId]!.quantiry > 1) {
      _items!.update(
        productId,
        (value) => CartItem(
          id: value.id,
          title: value.title,
          quantiry: value.quantiry - 1,
          price: value.price,
        ),
      );
    } else {
      _items!.remove(productId);
    }
  }

  void clear() {
    _items = {};
    notifyListeners();
  }
}
