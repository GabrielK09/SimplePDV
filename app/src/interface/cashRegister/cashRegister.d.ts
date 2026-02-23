interface CashRegisterContract {
    readonly id: number;
    readonly sale_id: number;
    readonly shopping_id: number;
    description: string;
    customer: string;
    specie_id: number;
    specie: string;
    input_value: number;
    output_value: number;
    total_balance: number;
};
