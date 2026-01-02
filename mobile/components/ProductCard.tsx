import { View, Text, Image, TouchableOpacity } from 'react-native';
import { styled } from 'nativewind';
import { Plus } from 'lucide-react-native';
import { Link } from 'expo-router';

const StyledView = styled(View);
const StyledText = styled(Text);
const StyledImage = styled(Image);

interface Product {
    id: string;
    name: string;
    price: string;
    image: string; // URL
    weight: string;
    rating: number;
}

interface ProductCardProps {
    product: Product;
}

export const ProductCard = ({ product }: ProductCardProps) => {
    return (
        <Link href={`/product/${product.id}`} asChild>
            <TouchableOpacity>
                <StyledView className="bg-card-light dark:bg-card-dark rounded-2xl p-3 shadow-sm mr-4 w-40">
                    {/* Image Container */}
                    <View className="h-32 w-full bg-gray-100 rounded-xl mb-3 overflow-hidden">
                        <StyledImage
                            source={{ uri: product.image }}
                            className="w-full h-full"
                            resizeMode="cover"
                        />
                    </View>

                    {/* Content */}
                    <StyledText className="font-bold text-gray-900 dark:text-white text-md mb-1" numberOfLines={1}>
                        {product.name}
                    </StyledText>

                    <View className="flex-row items-center space-x-1 mb-2">
                        {/* Rating Star Stub */}
                        <StyledText className="text-accent text-xs">★ {product.rating}</StyledText>
                        <StyledText className="text-gray-400 text-xs">({product.weight})</StyledText>
                    </View>

                    {/* Price and Action Row */}
                    <View className="flex-row justify-between items-center">
                        <StyledText className="font-bold text-lg text-gray-900 dark:text-white">
                            {product.price}
                        </StyledText>

                        {/* Add to Cart Button - Ensuring high contrast in dark mode as requested */}
                        <TouchableOpacity className="bg-primary p-2 rounded-full">
                            <Plus size={20} color="#FFFFFF" />
                        </TouchableOpacity>
                    </View>
                </StyledView>
            </TouchableOpacity>
        </Link>
    );
};
