interface ShoppingContract {
    id: number|null;
    load: number|null;
    shopping_itens: ShoppingItemContract[];
    total_shopping: number|null;
    status?: 'Pendente'|'Cancelado'|'Concluída'|'';
};

interface ShoppingItemContract {
    product_id: number|null;
    name: string|null;
    purchased_value: number|null;
    qtde_purchased: number|null;
    productWithCharacteristics?: ProductCharacteristicsContract[]|null;
};
