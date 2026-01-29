interface ShoppingContract {
    readonly id: number;
    load: number;
    operation: string;
    shoppingWithItem: ShoppingItem[]
};

interface ShoppingItem {
    readonly shoppingId: number;
    name: number;
    price: number;
    qtde: number;
};
