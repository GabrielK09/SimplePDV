export default function formatGridDataForPush(gridData: ProductCharacteristicsContract): ProductCharacteristicsContract {
    return {
        grid_qtde: Number(gridData.grid_qtde) || 0,
        id: gridData.id ?? null,
        product_id: gridData.product_id ?? null,
        size: gridData.size
    };
};
