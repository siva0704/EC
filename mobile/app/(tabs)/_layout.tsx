import { Tabs } from 'expo-router';
import { Home, Search, ShoppingCart, User } from 'lucide-react-native';

export default function TabLayout() {
    return (
        <Tabs screenOptions={{
            tabBarActiveTintColor: '#4ADE80',
            tabBarInactiveTintColor: '#9CA3AF',
            headerShown: false,
            tabBarStyle: {
                borderTopWidth: 0,
                elevation: 10,
                shadowOpacity: 0.1,
                height: 60,
                paddingBottom: 10,
            }
        }}>
            <Tabs.Screen
                name="home"
                options={{
                    title: 'Home',
                    tabBarIcon: ({ color }) => <Home size={24} color={color} />,
                }}
            />
            <Tabs.Screen
                name="search"
                options={{
                    title: 'Search',
                    tabBarIcon: ({ color }) => <Search size={24} color={color} />,
                }}
            />
            {/* Stubs for other tabs */}
        </Tabs>
    );
}
