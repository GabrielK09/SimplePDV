interface ProductContract {
    readonly id: number;
    name: string;
    price: number|null;
    qtde: number|null;
    commission: number|null;
    use_grid?: boolean;
    productWithCharacteristics?: ProductCharacteristicsContract[]|null
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
    readonly id: number;
    product_id: number;
    grid_qtde: number;
    size: Sizes|null;
};