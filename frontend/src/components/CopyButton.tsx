import { Button, message } from 'antd';
import { CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface Props {
    text: string
}
export const CopyButton = (props: Props) => {
    const { text } = props;
      const { t } = useTranslation();
    
    const handleCopy = () => {
        navigator.clipboard.writeText(text)
            .then(() => {
                message.success(t('content.copiedSuccessfully'));
            })
            .catch(() => {
                message.error(t('content.copyFailed'));
            });
    };

    return (
        <Button
            icon={<CopyOutlined />}
            onClick={handleCopy}
        >
            {t('content.copy')}
        </Button>
    );
};
