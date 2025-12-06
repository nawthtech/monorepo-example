import React from 'react';
import { Box, Typography, Container } from '@mui/material';

const MediaStudio: React.FC = () => {
  return (
    <Container maxWidth="lg">
      <Box sx={{ py: 4 }}>
        <Typography variant="h4" gutterBottom>
          🎨 استوديو الوسائط
        </Typography>
        <Typography variant="body1" color="text.secondary">
          تصميم وتحرير الوسائط باستخدام الذكاء الاصطناعي
        </Typography>
        
        <Box sx={{ mt: 4, p: 3, bgcolor: 'background.paper', borderRadius: 2 }}>
          <Typography variant="h6" gutterBottom>
            🚧 قيد الإنشاء
          </Typography>
          <Typography>
            استوديو الوسائط المتكامل قيد التطوير حالياً.
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default MediaStudio;