import { TouchableOpacity, Text, View } from 'react-native';
import { styled } from 'nativewind';

const StyledButton = styled(TouchableOpacity);
const StyledText = styled(Text);

interface PrimaryButtonProps {
    title: string;
    onPress: () => void;
    variant?: 'primary' | 'secondary' | 'outline';
    disabled?: boolean;
}

export const PrimaryButton = ({ title, onPress, variant = 'primary', disabled }: PrimaryButtonProps) => {
    const bgColors = {
        primary: 'bg-primary',
        secondary: 'bg-secondary',
        outline: 'bg-transparent border-2 border-primary',
    };

    const textColors = {
        primary: 'text-white',
        secondary: 'text-white',
        outline: 'text-primary',
    };

    return (
        <StyledButton
            onPress={onPress}
            disabled={disabled}
            className={`${bgColors[variant]} p-4 rounded-xl items-center justify-center w-full shadow-sm active:opacity-80`}
        >
            <StyledText className={`${textColors[variant]} font-bold text-lg`}>
                {title}
            </StyledText>
        </StyledButton>
    );
};
