# Turkcell Recommendation System - Frontend

A modern React/Next.js frontend for the Turkcell package recommendation system. This application provides an intuitive wizard-based interface for customers to get personalized mobile, home internet, and TV package recommendations.

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ 
- npm, yarn, or pnpm
- Backend API running on `http://localhost:8000`

### Installation & Setup

1. **Install dependencies:**
```bash
npm install
# or
yarn install
# or
pnpm install
```

2. **Environment Configuration:**
Create a `.env.local` file in the frontend directory:
```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8000

# Optional: For production builds
NODE_ENV=production
```

3. **Start development server:**
```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

4. **Open your browser:**
Navigate to [http://localhost:3000](http://localhost:3000)

## 🛠️ Development Scripts

```bash
# Development server with Turbopack
npm run dev

# Production build
npm run build

# Start production server
npm run start

# Run ESLint
npm run lint
```

## 🏗️ Project Structure

```
frontend/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── page.tsx           # Home/Setup wizard
│   │   ├── recommendations/   # Recommendations page
│   │   ├── checkout/          # Checkout flow
│   │   ├── layout.tsx         # Root layout with providers
│   │   └── globals.css        # Global styles
│   ├── components/            # Reusable UI components
│   │   ├── HouseholdForm.tsx  # Dynamic household input
│   │   ├── AddressForm.tsx    # Address & coverage input
│   │   ├── CoverageBadge.tsx  # Technology coverage display
│   │   ├── RecommendationCard.tsx # Package recommendation cards
│   │   ├── SummaryModal.tsx   # Detailed cost breakdown
│   │   └── SlotPicker.tsx     # Installation slot selection
│   ├── context/               # React Context providers
│   │   └── WizardContext.tsx  # Wizard state management
│   ├── lib/                   # Utilities and configurations
│   │   ├── api-client.ts      # API communication layer
│   │   ├── hooks.ts           # React Query data hooks
│   │   └── query-client.tsx   # React Query provider setup
│   └── types/                 # TypeScript type definitions
│       └── api.ts             # API request/response types
├── public/                    # Static assets
├── package.json               # Dependencies and scripts
└── README.md                  # This file
```

## 🎯 Features

### Multi-Step Wizard Flow
1. **Setup**: Household information and address input
2. **Recommendations**: AI-powered package suggestions
3. **Checkout**: Installation scheduling and order confirmation

### Core Components

#### 🏠 HouseholdForm
- Dynamic line management (add/remove household members)
- Usage input: expected GB, minutes, TV hours
- Real-time validation and visual feedback
- Usage summary with totals

#### 📍 AddressForm  
- Address selection with sample addresses
- Real-time coverage lookup
- Technology availability badges
- Coverage information display

#### 💳 RecommendationCard
- Beautiful package presentation
- Cost breakdown with savings
- Service details (mobile/home/TV)
- Rank-based styling and recommendations

#### 📅 SlotPicker
- Calendar-style installation slot selection
- Technology-specific scheduling
- Duration and availability indicators
- Grouped by date with time slots

#### 📊 SummaryModal
- Detailed cost breakdown
- Individual line costs with overage
- Discount explanations
- Professional invoice-style layout

### State Management
- **React Context**: Wizard state across all steps
- **React Query**: Server state and caching
- **TypeScript**: Full type safety throughout

## 🔧 Technology Stack

### Core Framework
- **Next.js 15**: React framework with App Router
- **React 19**: Latest React with concurrent features  
- **TypeScript**: Full type safety and developer experience

### Styling & UI
- **Tailwind CSS 4**: Utility-first CSS framework
- **Custom Components**: Reusable UI component library
- **Responsive Design**: Mobile-first responsive layouts

### Data Management
- **React Query v5**: Server state management and caching
- **React Context**: Client-side state management
- **Custom Hooks**: Reusable data fetching logic

### Developer Experience
- **ESLint**: Code linting and style enforcement
- **TypeScript**: Static type checking
- **Turbopack**: Fast development builds

## 🌐 API Integration

### Endpoints Used
- `GET /health` - API health monitoring
- `GET /api/coverage/{address_id}` - Technology coverage
- `GET /api/install-slots/{address_id}` - Installation scheduling
- `POST /api/recommendation` - Package recommendations
- `POST /api/checkout` - Order processing

### Data Hooks
```typescript
// Coverage and installation slots
const { coverage, slots, isLoading } = useCatalog(addressId);

// Package recommendations
const recommendationMutation = useRecommendation();
recommendationMutation.mutate(requestData);

// Order processing
const checkoutMutation = useCheckout();
checkoutMutation.mutate(checkoutData);
```

## 🎨 Design System

### Color Palette
- **Primary**: Blue tones for CTAs and navigation
- **Success**: Green for confirmations and positive states
- **Warning**: Yellow for information and alerts
- **Error**: Red for errors and critical states
- **Neutral**: Gray scale for text and backgrounds

### Typography
- **Headings**: Bold weights for hierarchy
- **Body**: Regular weights for readability
- **Code**: Monospace for technical content

### Components
- **Cards**: Elevated surfaces with shadows
- **Buttons**: Consistent styling across states
- **Forms**: Clear labels and validation
- **Modals**: Focused overlays for details

## 🔍 User Flow

### 1. Setup Wizard
1. User enters household information (lines, usage patterns)
2. User selects address and sees coverage options
3. Form validation ensures data completeness
4. Navigation to recommendations

### 2. Recommendations
1. System analyzes requirements and returns top 3 packages
2. User reviews options with cost breakdowns
3. Detailed modals show complete cost analysis
4. User selects preferred package

### 3. Checkout
1. Order summary displays selected package
2. User selects installation time slot
3. Order confirmation with generated order ID
4. Success page with next steps

## 🧪 Testing & Development

### Manual Testing
1. **Start the backend API** (see backend README)
2. **Run frontend development server**: `npm run dev`
3. **Navigate through the wizard**:
   - Add household lines with various usage patterns
   - Select different addresses (A1001-A1005 for testing)
   - Review recommendations and detailed breakdowns
   - Complete checkout flow with slot selection

### Sample Test Data
```json
{
  "household": [
    {
      "line_id": "LINE001",
      "expected_gb": 8.0,
      "expected_min": 450.0,
      "tv_hd_hours": 25.0
    }
  ],
  "address_id": "A1001"
}
```

### Browser Testing
- **Chrome/Edge**: Primary testing browsers
- **Firefox**: Secondary browser support
- **Safari**: Mobile and desktop compatibility
- **Mobile**: Responsive design testing

## 🚀 Production Deployment

### Build Process
```bash
# Create production build
npm run build

# Start production server
npm run start
```

### Environment Variables
```bash
# Production API endpoint
NEXT_PUBLIC_API_URL=https://api.turkcell-recommendations.com

# Optional build settings
NODE_ENV=production
NEXT_TELEMETRY_DISABLED=1
```

### Performance Optimizations
- **Static Generation**: Pre-built pages where possible
- **Image Optimization**: Next.js automatic image optimization
- **Code Splitting**: Automatic route-based code splitting
- **Caching**: React Query automatic request caching

## 🐛 Troubleshooting

### Common Issues

#### "API connection failed"
- Ensure backend server is running on port 8000
- Check CORS configuration in backend
- Verify API endpoint in environment variables

#### "Coverage information not loading"
- Verify address ID format (A1001, A1002, etc.)
- Check backend database connectivity
- Ensure sample data is seeded

#### "Recommendations not appearing"
- Validate household form completion
- Check network requests in browser dev tools
- Verify API request format matches backend expectations

#### "Checkout fails"
- Ensure installation slot is selected
- Verify all required fields are present
- Check backend logs for validation errors

### Debug Mode
Enable React Query DevTools in development:
- Open browser dev tools
- Look for React Query tab
- Monitor cache state and network requests

## 📱 Responsive Design

### Breakpoints
- **Mobile**: < 640px (sm)
- **Tablet**: 640px - 1024px (md-lg)
- **Desktop**: > 1024px (xl+)

### Mobile Optimizations
- Touch-friendly button sizes
- Optimized form layouts
- Simplified navigation
- Fast loading times

## 🔒 Security

### Data Handling
- No sensitive data stored in client
- Secure API communication
- Input validation and sanitization
- XSS protection via React

### Best Practices
- Environment variable usage
- Secure HTTP headers
- Content Security Policy
- Regular dependency updates

---

**Frontend Version**: 0.1.0  
**Framework**: Next.js 15 + React 19  
**Last Updated**: December 2024  
**Contact**: Turkcell Frontend Team

## 🚀 Getting Help

1. **Check this README** for common setup issues
2. **Review browser console** for client-side errors  
3. **Inspect Network tab** for API communication issues
4. **Check backend logs** for server-side problems
5. **Verify environment variables** are correctly set