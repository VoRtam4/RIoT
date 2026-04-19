import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";

type Props = {
  instance: {
    id: string;
    label?: string | null;
    uid?: string | null;
  };
  selected: boolean;
  onClick: () => void;
};

export default function KpiInstanceCard({
  instance,
  selected,
  onClick,
}: Props) {
  return (
    <Card
      sx={{
        flexShrink: 0,
        height: "100%",
        overflow: "hidden",
        backgroundColor: selected
          ? "var(--primary)"
          : "var(--bg-card)",
        transition: "all 0.15s",
        border: selected
          ? "1px solid var(--bs-primary)"
          : "1px solid transparent",
        "&:hover": {
          transform: "translateY(-1px)",
          boxShadow: 3,
        },
      }}
    >
      <CardActionArea
        onClick={onClick}
        sx={{
          height: "100%",
          backgroundColor: selected
            ? "rgba(0,123,255,0.1)"
            : "transparent",
          "&:hover": {
            backgroundColor: "rgba(0,123,255,0.2)",
          },
        }}
      >
        <CardContent
          sx={{
            py: 1.5,
            height: "100%",
            display: "flex",
            alignItems: "center",
            overflow: "hidden",
          }}
        >
          <Typography
            variant="body1"
            sx={{
              display: "-webkit-box",
              overflow: "hidden",
              textOverflow: "ellipsis",
              WebkitBoxOrient: "vertical",
              WebkitLineClamp: 2,
            }}
          >
            {instance.label || instance.uid || instance.id}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
