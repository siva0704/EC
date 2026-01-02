import { View, Text, Image, ScrollView, TouchableOpacity } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { styled } from 'nativewind';
import { useProducts } from '../../hooks/useProducts';
import { ArrowLeft, Heart } from 'lucide-react-native';

const StyledText = styled(Text);

export default function ProductDetailScreen() {
    const { id } = useLocalSearchParams();
    const router = useRouter();
    const { products } = useProducts();
    const product = products.find(p => p.id === id) || products[0];

    return (
        <View className="flex-1 bg-white">
            <ScrollView>
                {/* Parallax Header Placeholder */}
                <View className="h-72 bg-gray-100 relative">
                    <Image source={{ uri: product.image }} className="w-full h-full" resizeMode="cover" />
                    <TouchableOpacity onPress={() => router.back()} className="absolute top-12 left-4 bg-white p-2 rounded-full">
                        <ArrowLeft size={24} color="black" />
                    </TouchableOpacity>
                </View>

                <View className="p-6 -mt-6 bg-white rounded-t-3xl shadow-lg h-full">
                    <View className="flex-row justify-between items-start mb-4">
                        <View>
                            <StyledText className="text-2xl font-bold text-gray-900">{product.name}</StyledText>
                            <StyledText className="text-lg text-primary font-bold">{product.price}</StyledText>
                        </View>
                        <TouchableOpacity className="bg-red-50 p-2 rounded-full">
                            <Heart size={24} color="#EF4444" />
                        </TouchableOpacity>
                    </View>

                    <StyledText className="text-gray-600 leading-6 mb-6">
                        {product.description}
                    </StyledText>

                    {/* Weight Selector Stub */}
                    <StyledText className="font-bold mb-3">Weight</StyledText>
                    <View className="flex-row gap-3 mb-8">
                        <View className="bg-primary px-4 py-2 rounded-full"><Text className="text-white font-bold">1kg</Text></View>
                        <View className="bg-gray-100 px-4 py-2 rounded-full"><Text className="text-gray-600">5kg</Text></View>
                    </View>
                </View>
            </ScrollView>

            {/* Sticky Footer */}
            <View className="absolute bottom-0 w-full bg-white p-4 border-t border-gray-100 flex-row gap-4 safe-area-bottom">
                <TouchableOpacity className="flex-1 bg-secondary py-4 rounded-xl items-center">
                    <Text className="text-white font-bold text-lg">Add to Cart</Text>
                </TouchableOpacity>
                <TouchableOpacity
                    className="flex-1 bg-primary py-4 rounded-xl items-center"
                    onPress={() => router.push('/checkout/success')}
                >
                    <Text className="text-white font-bold text-lg">Buy Now</Text>
                </TouchableOpacity>
            </View>
        </View>
    );
}
