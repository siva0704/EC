export const useProducts = () => {
    const products = [
        {
            id: "1",
            name: "Royal Basmati Rice",
            price: "$8.99",
            weight: "10lb",
            rating: 4.8,
            image: "https://example.com/basmati.png", // Placeholder
            description: "Premium aged basmati rice with long grains and aromatic flavor."
        },
        {
            id: "2",
            name: "Nishiki Jasmine",
            price: "$12.50",
            weight: "5lb",
            rating: 4.9,
            image: "https://example.com/jasmine.png",
            description: "High quality jasmine rice, perfect for sushi and sticky rice dishes."
        },
        {
            id: "3",
            name: "Organic Brown Rice",
            price: "$15.50",
            weight: "5kg",
            rating: 4.7,
            image: "https://example.com/brown.png",
            description: "Whole grain brown rice, rich in fiber and nutrients."
        }
    ];

    return { products };
};
