interface ShoppingContract {
    readonly id: number;
    load: number;
    shoppingWithItem: ProductContract[]
    totalShopping: number;
};
