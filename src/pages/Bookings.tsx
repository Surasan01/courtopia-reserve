import { useState, useEffect } from "react";
import { format } from "date-fns";
import Header from "@/components/Header";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useToast } from "@/hooks/use-toast";
import { fetchUserBookings, cancelBooking } from "@/lib/api";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Calendar, Clock } from "lucide-react";

interface Booking {
  id: string;
  courtNumber: number;
  courtName: string;
  date: string;
  time: string;
  status: string;
}

const Bookings = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // แยกการจองตามสถานะ
  const activeBookings = bookings.filter(booking => booking.status === 'active');
  const completedBookings = bookings.filter(booking => booking.status === 'completed');
  const cancelledBookings = bookings.filter(booking => booking.status === 'cancelled');

  useEffect(() => {
    const loadUserBookings = async () => {
      if (!user) {
        navigate("/login");
        return;
      }

      setIsLoading(true);

      try {
        const response = await fetchUserBookings(user.token);

        if (!response || !Array.isArray(response) || response.length === 0) {
          setBookings([]);
          return;
        }

        // Transform the API response to match our Booking interface
        const transformedBookings = response.map((booking: any) => ({
          id: booking.id,
          courtNumber: booking.courtNumber,
          courtName: `คอร์ท ${booking.courtNumber}`,
          date: booking.bookingDate,
          time: `${booking.startTime} - ${booking.endTime}`,
          status: booking.status
        }));

        setBookings(transformedBookings);
      } catch (error) {
        toast({
          title: 'ไม่สามารถโหลดข้อมูลการจองได้',
          description: error instanceof Error ? error.message : 'กรุณาลองใหม่อีกครั้ง',
          variant: 'destructive',
        });
        
        // ไม่ใช้ข้อมูลตัวอย่างในโปรดักชันจริง
        setBookings([]);
      } finally {
        setIsLoading(false);
      }
    };

    loadUserBookings();
  }, [navigate, user, toast]);

  const handleCancelBooking = async (id: string) => {
    if (!user) return;

    try {
      await cancelBooking(user.token, id);

      // Update the booking status instead of removing it
      setBookings(bookings.map(booking => 
        booking.id === id 
          ? { ...booking, status: 'cancelled' } 
          : booking
      ));

      toast({
        title: "ยกเลิกการจองสำเร็จ",
        description: "คุณได้ยกเลิกการจองเรียบร้อยแล้ว",
      });
    } catch (error) {
      toast({
        title: "ไม่สามารถยกเลิกการจองได้",
        description: error instanceof Error ? error.message : "กรุณาลองใหม่อีกครั้ง",
        variant: "destructive",
      });
    }
  };

  // Component สำหรับแสดงการจอง
  const BookingCard = ({ booking, showCancelButton = false }: { booking: Booking, showCancelButton?: boolean }) => (
    <div key={booking.id} className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 flex flex-wrap md:flex-nowrap justify-between items-center">
      <div className="mb-4 md:mb-0 flex-1">
        <div className="flex items-center justify-between mb-2">
          <h3 className="font-semibold text-lg">{booking.courtName}</h3>
          <StatusBadge status={booking.status} />
        </div>
        <div className="flex items-center gap-2 text-gray-600">
          <Calendar size={16} />
          <p className="text-sm">{booking.date}</p>
        </div>
        <div className="flex items-center gap-2 text-gray-600">
          <Clock size={16} />
          <p className="text-sm">{booking.time}</p>
        </div>
      </div>
      {showCancelButton && (
        <Button
          variant="destructive"
          size="sm"
          onClick={() => handleCancelBooking(booking.id)}
        >
          ยกเลิกการจอง
        </Button>
      )}
    </div>
  );

  // Component สำหรับแสดงสถานะ
  const StatusBadge = ({ status }: { status: string }) => {
    switch (status) {
      case 'active':
        return <Badge className="bg-green-500">จองแล้ว</Badge>;
      case 'completed':
        return <Badge className="bg-blue-500">จบแล้ว</Badge>;
      case 'cancelled':
        return <Badge className="bg-red-500">ยกเลิก</Badge>;
      default:
        return <Badge>{status}</Badge>;
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <main className="flex-1 bg-court-white">
        <div className="container mx-auto px-4 py-8">
          <h1 className="text-2xl font-bold mb-6">การจองของฉัน</h1>

          {isLoading ? (
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-court-orange"></div>
            </div>
          ) : bookings.length > 0 ? (
            <Tabs defaultValue="active">
              <TabsList className="mb-4">
                <TabsTrigger value="active">
                  กำลังจองอยู่ 
                  {activeBookings.length > 0 && (
                    <span className="ml-2 inline-flex items-center justify-center w-5 h-5 text-xs font-semibold rounded-full bg-court-orange text-white">
                      {activeBookings.length}
                    </span>
                  )}
                </TabsTrigger>
                <TabsTrigger value="completed">
                  จบแล้ว
                  {completedBookings.length > 0 && (
                    <span className="ml-2 inline-flex items-center justify-center w-5 h-5 text-xs font-semibold rounded-full bg-blue-500 text-white">
                      {completedBookings.length}
                    </span>
                  )}
                </TabsTrigger>
                <TabsTrigger value="cancelled">
                  ยกเลิกแล้ว
                  {cancelledBookings.length > 0 && (
                    <span className="ml-2 inline-flex items-center justify-center w-5 h-5 text-xs font-semibold rounded-full bg-red-500 text-white">
                      {cancelledBookings.length}
                    </span>
                  )}
                </TabsTrigger>
              </TabsList>
              
              <TabsContent value="active" className="space-y-4">
                {activeBookings.length > 0 ? (
                  activeBookings.map(booking => (
                    <BookingCard key={booking.id} booking={booking} showCancelButton={true} />
                  ))
                ) : (
                  <div className="text-center py-8 bg-gray-50 rounded-lg">
                    <p className="text-gray-500">ไม่มีการจองที่กำลังดำเนินอยู่</p>
                  </div>
                )}
              </TabsContent>
              
              <TabsContent value="completed" className="space-y-4">
                {completedBookings.length > 0 ? (
                  completedBookings.map(booking => (
                    <BookingCard key={booking.id} booking={booking} />
                  ))
                ) : (
                  <div className="text-center py-8 bg-gray-50 rounded-lg">
                    <p className="text-gray-500">ไม่มีประวัติการจองที่สำเร็จ</p>
                  </div>
                )}
              </TabsContent>
              
              <TabsContent value="cancelled" className="space-y-4">
                {cancelledBookings.length > 0 ? (
                  cancelledBookings.map(booking => (
                    <BookingCard key={booking.id} booking={booking} />
                  ))
                ) : (
                  <div className="text-center py-8 bg-gray-50 rounded-lg">
                    <p className="text-gray-500">ไม่มีการจองที่ถูกยกเลิก</p>
                  </div>
                )}
              </TabsContent>
            </Tabs>
          ) : (
            <div className="text-center py-12 bg-gray-50 rounded-lg">
              <p className="text-gray-500 mb-4">ไม่มีการจองในขณะนี้</p>
              <Button
                className="bg-court-orange hover:bg-court-orange/90"
                onClick={() => navigate("/")}
              >
                จองคอร์ทเลย
              </Button>
            </div>
          )}
        </div>
      </main>

      <footer className="bg-court-orange text-white p-4 text-center">
        <p>© {new Date().getFullYear()} Courtminton - King Mongkut's University of Technology Thonburi</p>
      </footer>
    </div>
  );
};

export default Bookings;