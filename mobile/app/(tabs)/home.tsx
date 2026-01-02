import { View, Text, FlatList, TextInput, ScrollView, SafeAreaView } from 'react-native';
import { styled } from 'nativewind';
import { ProductCard } from '../../components/ProductCard';
import { useProducts } from '../../hooks/useProducts';
import { Menu, Search } from 'lucide-react-native';

const StyledView = styled(View);
const StyledText = styled(Text);

export default function HomeScreen() {
    const { products } = useProducts();

    const categories = [
        { id: '1', name: 'Basmati', color: '#FFFHT' },
        { id: '2', name: 'Jasmine', color: '#FFFHT' },
        { id: '3', name: 'Brown', color: '#FFFHT' },
        { id: '4', name: 'Arborio', color: '#FFFHT' },
    ];

    return (
        <SafeAreaView className="flex-1 bg-surface-light dark:bg-black">
            <ScrollView className="p-4" showsVerticalScrollIndicator={false}>
                {/* Header */}
                <View className="flex-row items-center justify-between mb-6">
                    <Menu color="#374151" size={24} />
                    <View className="flex-1 mx-4 bg-white rounded-full flex-row items-center px-4 py-2 shadow-sm">
                        <Search size={20} color="#9CA3AF" />
                        <TextInput placeholder="Search for rice..." className="ml-2 flex-1" />
                    </View>
                </View>

                {/* Categories */}
                <StyledText className="text-xl font-bold text-gray-900 mb-4">Categories</StyledText>
                <FlatList
                    horizontal
                    data={categories}
                    keyExtractor={item => item.id}
                    showsHorizontalScrollIndicator={false}
                    className="mb-8"
                    renderItem={({ item }) => (
                        <View className="mr-6 items-center">
                            <View className="w-16 h-16 rounded-full bg-white shadow-sm mb-2 items-center justify-center">
                                {/* Placeholder for Rice Image */}
                                <View className="w-12 h-12 bg-gray-200 rounded-full" />
                            </View>
                            <StyledText className="text-sm font-medium text-gray-700">{item.name}</StyledText>
                        </View>
                    )}
                />

                {/* Featured */}
                <StyledText className="text-xl font-bold text-gray-900 mb-4">Featured Products</StyledText>
                <FlatList
                    horizontal
                    data={products}
                    keyExtractor={item => item.id}
                    showsHorizontalScrollIndicator={false}
                    renderItem={({ item }) => <ProductCard product={item} />}
                />
            </ScrollView>
        </SafeAreaView>
    );
}
