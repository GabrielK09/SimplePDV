interface DashBoardContract {
    commission: number;
	total_saled: number;
	best_customer: string;
	amount_saled: number;
    total_shopping: number;
    amount_shopping_itens: number;
    amount_shopping: number;
};

interface PopularItensFilterContract {
    per_page: number|null;
};
