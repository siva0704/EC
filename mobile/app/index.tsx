import { View, Text, Image } from 'react-native';
import { useRouter } from 'expo-router';
import { PrimaryButton } from '../components/PrimaryButton';
import { styled } from 'nativewind';

const StyledView = styled(View);
const StyledText = styled(Text);

export default function OnboardingScreen() {
    const router = useRouter();

    return (
        <StyledView className="flex-1 bg-white items-center justify-between p-6 pt-20 pb-10">
            <View className="items-center">
                {/* Illustration Stub */}
                <View className="w-64 h-64 bg-surface-light rounded-full mb-8 items-center justify-center">
                    <StyledText className="text-4xl">🌾</StyledText>
                </View>

                <StyledText className="text-3xl font-bold text-center text-gray-900 mb-4">
                    Fresh Grains, Delivered to You
                </StyledText>
                <StyledText className="text-lg text-center text-gray-500 px-4">
                    Explore a world of premium rice varieties and get them delivered fresh to your doorstep.
                </StyledText>
            </View>

            <View className="w-full">
                <PrimaryButton
                    title="Get Started"
                    onPress={() => router.replace('/(tabs)/home')}
                />
            </View>
        </StyledView>
    );
}
