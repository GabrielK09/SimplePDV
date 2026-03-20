interface ShoppingContract {
    id: number|null;
    load: number|null;
    shopping_itens: ShoppingItemContract[];
    total_shopping: number|null;
    status?: 'Pendente'|'Cancelado'|'Concluída'|'';
};

interface ShoppingItemContract {
    readonly product_id: number;
    name: string;
    purchased_value: number;
    qtde_purchased: number;
};
