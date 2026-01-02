import { View, Text } from 'react-native';
import { PrimaryButton } from '../../components/PrimaryButton';
import { useRouter } from 'expo-router';
import { CheckCircle } from 'lucide-react-native';

export default function SuccessScreen() {
    const router = useRouter();

    return (
        <View className="flex-1 bg-surface-light items-center justify-center p-6">
            <CheckCircle size={80} color="#4ADE80" className="mb-6" />
            <Text className="text-3xl font-bold text-gray-900 mb-2">Order Confirmed!</Text>
            <Text className="text-gray-500 text-center mb-8">
                Your rice is on its way. Estimated delivery: 24 hours.
            </Text>

            <PrimaryButton title="Back Home" onPress={() => router.navigate('/(tabs)/home')} />
        </View>
    );
}
