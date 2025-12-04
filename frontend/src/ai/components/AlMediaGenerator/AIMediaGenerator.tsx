import React, { useState } from 'react';
import { useAI } from '../../hooks/useAI';
import {
  Box,
  Button,
  TextField,
  Card,
  CardContent,
  Typography,
  Grid,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  CircularProgress,
  Alert,
  ImageList,
  ImageListItem,
  IconButton,
  Chip,
} from '@mui/material';
import {
  Image,
  PlayCircle,
  Download,
  Share,
  Refresh,
} from '@mui/icons-material';

const AIMediaGenerator: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [mediaType, setMediaType] = useState<'image' | 'video'>('image');
  const [style, setStyle] = useState('realistic');
  const [generatedMedia, setGeneratedMedia] = useState<Array<{
    id: string;
    url: string;
    type: string;
    prompt: string;
    timestamp: Date;
  }>>([]);
  
  const { generateImage, loading, error } = useAI();
  
  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    
    try {
      const result = await generateImage(prompt, style);
      
      if (result.success) {
        const newMedia = {
          id: Date.now().toString(),
          url: result.data.url,
          type: mediaType,
          prompt,
          timestamp: new Date(),
        };
        
        setGeneratedMedia(prev => [newMedia, ...prev.slice(0, 9)]);
      }
    } catch (err) {
      console.error('Generation failed:', err);
    }
  };
  
  const downloadMedia = (url: string, filename: string) => {
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
  };
  
  return (
    <Box sx={{ maxWidth: 1200, margin: '0 auto', p: 3 }}>
      <Typography variant="h4" gutterBottom sx={{ color: '#7A3EF0', mb: 4 }}>
        🎨 NawthTech AI Media Generator
      </Typography>
      
      <Grid container spacing={3}>
        {/* لوحة التحكم */}
        <Grid item xs={12} md={4}>
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                إعدادات الوسائط
              </Typography>
              
              <FormControl fullWidth sx={{ mb: 2 }}>
                <InputLabel>نوع الوسائط</InputLabel>
                <Select
                  value={mediaType}
                  label="نوع الوسائط"
                  onChange={(e) => setMediaType(e.target.value as any)}
                >
                  <MenuItem value="image">صورة</MenuItem>
                  <MenuItem value="video" disabled>فيديو (قريباً)</MenuItem>
                </Select>
              </FormControl>
              
              <TextField
                fullWidth
                multiline
                rows={3}
                label="وصف الصورة"
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                sx={{ mb: 2 }}
                placeholder="مثال: شعار عصري لشركة تقنية بالألوان الأرجواني والسماوي"
              />
              
              <FormControl fullWidth sx={{ mb: 3 }}>
                <InputLabel>النمط</InputLabel>
                <Select
                  value={style}
                  label="النمط"
                  onChange={(e) => setStyle(e.target.value)}
                >
                  <MenuItem value="realistic">واقعي</MenuItem>
                  <MenuItem value="anime">أنمي</MenuItem>
                  <MenuItem value="digital_art">فن رقمي</MenuItem>
                  <MenuItem value="3d_render">تصميم ثلاثي الأبعاد</MenuItem>
                  <MenuItem value="minimalist">بسيط</MenuItem>
                </Select>
              </FormControl>
              
              <Button
                fullWidth
                variant="contained"
                onClick={handleGenerate}
                disabled={loading || !prompt.trim()}
                sx={{
                  bgcolor: '#7A3EF0',
                  '&:hover': { bgcolor: '#6A2EE0' },
                  mb: 2,
                }}
                startIcon={loading ? <CircularProgress size={20} color="inherit" /> : <Image />}
              >
                {loading ? 'جاري التوليد...' : 'توليد وسائط'}
              </Button>
              
              {error && (
                <Alert severity="error" sx={{ mt: 2 }}>
                  {error}
                </Alert>
              )}
              
              <Alert severity="info" sx={{ mt: 2 }}>
                نصائح: كن وصفياً في الطلب، أضف تفاصيل عن الألوان والأسلوب والمزاج
              </Alert>
            </CardContent>
          </Card>
          
          {/* الإحصائيات */}
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                إحصائيات الاستخدام
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    الصور المولدة
                  </Typography>
                  <Typography variant="h5">
                    {generatedMedia.length}
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    الحصة المتبقية
                  </Typography>
                  <Typography variant="h5" color="success.main">
                    8/10
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
        
        {/* المعرض */}
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                <Typography variant="h6">
                  المعرض
                </Typography>
                <Chip 
                  label={`${generatedMedia.length} عنصر`} 
                  color="primary" 
                  variant="outlined" 
                />
              </Box>
              
              {generatedMedia.length > 0 ? (
                <ImageList cols={2} gap={16}>
                  {generatedMedia.map((item) => (
                    <ImageListItem 
                      key={item.id}
                      sx={{
                        borderRadius: 2,
                        overflow: 'hidden',
                        position: 'relative',
                        '&:hover .media-actions': {
                          opacity: 1,
                        },
                      }}
                    >
                      <img
                        src={item.url}
                        alt={item.prompt}
                        loading="lazy"
                        style={{
                          width: '100%',
                          height: 200,
                          objectFit: 'cover',
                        }}
                      />
                      
                      {/* Overlay Actions */}
                      <Box
                        className="media-actions"
                        sx={{
                          position: 'absolute',
                          top: 0,
                          left: 0,
                          right: 0,
                          bottom: 0,
                          bgcolor: 'rgba(0,0,0,0.5)',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                          opacity: 0,
                          transition: 'opacity 0.3s',
                        }}
                      >
                        <IconButton
                          onClick={() => downloadMedia(item.url, `nawthtech_${item.id}.jpg`)}
                          sx={{ color: 'white', bgcolor: 'rgba(255,255,255,0.2)' }}
                        >
                          <Download />
                        </IconButton>
                        <IconButton
                          sx={{ color: 'white', bgcolor: 'rgba(255,255,255,0.2)', mx: 1 }}
                        >
                          <Share />
                        </IconButton>
                        <IconButton
                          onClick={() => setPrompt(item.prompt)}
                          sx={{ color: 'white', bgcolor: 'rgba(255,255,255,0.2)' }}
                        >
                          <Refresh />
                        </IconButton>
                      </Box>
                      
                      {/* Prompt Preview */}
                      <Box
                        sx={{
                          position: 'absolute',
                          bottom: 0,
                          left: 0,
                          right: 0,
                          bgcolor: 'rgba(0,0,0,0.7)',
                          p: 1,
                        }}
                      >
                        <Typography
                          variant="caption"
                          sx={{
                            color: 'white',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                            overflow: 'hidden',
                          }}
                        >
                          {item.prompt}
                        </Typography>
                      </Box>
                    </ImageListItem>
                  ))}
                </ImageList>
              ) : (
                <Box
                  sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    height: 400,
                    border: '2px dashed #ddd',
                    borderRadius: 1,
                    p: 3,
                  }}
                >
                  <Image sx={{ fontSize: 60, color: '#7A3EF0', mb: 2 }} />
                  <Typography variant="body1" color="text.secondary" align="center">
                    لم يتم توليد وسائط بعد
                  </Typography>
                  <Typography variant="body2" color="text.secondary" align="center" sx={{ mt: 1 }}>
                    اكتب وصفاً واضحاً واضغط على "توليد وسائط"
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default AIMediaGenerator;