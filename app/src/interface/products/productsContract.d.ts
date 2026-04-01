interface ProductContract {
    id: number|null;
    name: string;
    price: number|null;
    qtde: number|null;
    commission: number|null;
    use_grid?: boolean;
    productWithCharacteristics?: ProductCharacteristicsContract[]|null;
    deleted_at?: Date;
};

enum Sizes {
    PP = 'PP', 
    P = 'P', 
    M = 'M', 
    G = 'G', 
    GG = 'GG', 
    XG = 'XG', 
    XGG = 'XGG', 
    EG = 'EG', 
    EGG = 'EGG', 
    O = 'O'
};

interface ProductCharacteristicsContract {
    readonly id: number|null;
    product_id: number|null;
    grid_qtde: number|null;
    size: Sizes|null;
};