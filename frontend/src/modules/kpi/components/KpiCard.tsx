import { useNavigate } from "react-router-dom";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Chip from "@mui/material/Chip";

type Props = {
  kpi: {
    id: string;
    label?: string | null;
    sdTypeID: string;
    sdInstanceMode?: any;
  };
};
export default function KpiCard({ kpi }: Props) {
  const navigate = useNavigate();

  return (
    <Card
      sx={{
        height: "100%",
        minHeight: 120,
        backgroundColor: "var(--bg-card)",
        "&:hover": {
          backgroundColor: "var(--border)",
        },
      }}
    >
      <CardActionArea
        onClick={() => navigate(`/kpi/${kpi.id}`)}
        sx={{
          height: "100%",
          display: "flex",
          alignItems: "stretch",
        }}
      >
        <CardContent
          sx={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
            width: "100%",
          }}
        >
          <Typography variant="h6" className="color-label">
            {kpi.label ?? "\u00A0"}
          </Typography>

          <Box sx={{ mt: 1 }}>
            <Typography
              variant="body2"
              sx={{ opacity: 0.7 }}
              className="color-label"
            >
              {kpi.id ?? "\u00A0"}
            </Typography>

            <Box sx={{ mt: 1 }}>
              <Chip label={kpi.sdInstanceMode} size="small" color="primary" />
            </Box>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}