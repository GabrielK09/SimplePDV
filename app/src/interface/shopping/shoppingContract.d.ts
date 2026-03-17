interface ShoppingContract {
    readonly id: number;
    load: number;
    shopping_itens: ShoppingItemContract[];
    totalShopping: number;
};

interface ShoppingItemContract {
    readonly product_id: number;
    name: string;
    purchased_value: number;
    qtde_purchased: number;
};
