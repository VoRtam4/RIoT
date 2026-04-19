import { useNavigate } from "react-router-dom";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";

type Instance = {
  id: string;
  label?: string | null;
  uid?: string | null;
};

type Props = {
  instance: Instance;
};

export default function SdInstanceCard({ instance }: Props) {
  const navigate = useNavigate();
  return (
    <Card
      sx={{
        height: "100%",
        minHeight: 120,
        backgroundColor: "var(--bg-card)",
      }}
    >
      <CardActionArea
        onClick={() => navigate(`/sd-instance/${instance.id}`)}
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
          {/* TITLE */}
          <Typography variant="h6" className="color-label">
            {instance.label || instance.uid || instance.id}
          </Typography>

          {/* META */}
          <Box sx={{ mt: 1 }}>
            <Typography
              variant="body2"
              sx={{ opacity: 0.7 }}
              className="color-label"
            >
              {instance.uid ?? "\u00A0"}
            </Typography>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
