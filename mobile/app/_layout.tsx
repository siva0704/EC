import { Stack } from "expo-router";
import { View } from "react-native";

export default function RootLayout() {
    return (
        <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="index" options={{ headerShown: false }} />
            <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
            <Stack.Screen name="product/[id]" options={{ presentation: 'modal', headerShown: false }} />
            <Stack.Screen name="checkout/success" options={{ presentation: 'fullScreenModal' }} />
        </Stack>
    );
}
