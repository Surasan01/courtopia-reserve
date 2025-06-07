import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useToast } from '@/hooks/use-toast';
import { useAuth } from '@/hooks/useAuth';
import { login } from '@/lib/api';

import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Eye, EyeOff } from 'lucide-react';

const loginSchema = z.object({
    studentId: z.string().min(1, 'กรุณากรอกรหัสนักศึกษา'),
    password: z.string().min(1, 'กรุณากรอกรหัสผ่าน'),
});

type LoginFormValues = z.infer<typeof loginSchema>;

const Login = () => {
    const { login: authLogin } = useAuth();
    const navigate = useNavigate();
    const { toast } = useToast();
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const form = useForm<LoginFormValues>({
        resolver: zodResolver(loginSchema),
        defaultValues: {
            studentId: '',
            password: '',
        },
    });

    const onSubmit = async (data: LoginFormValues) => {
        setIsLoading(true);
        try {
            const response = await login(data.studentId, data.password);

            toast({
                title: 'เข้าสู่ระบบสำเร็จ',
                description: `ยินดีต้อนรับ ${response.name}`,
            });

            authLogin(response.token, response.studentId, response.name, response.role);

            navigate('/');
        } catch (error) {
            toast({
                title: 'เข้าสู่ระบบไม่สำเร็จ',
                description: error instanceof Error ? error.message : 'เกิดข้อผิดพลาด กรุณาลองอีกครั้ง',
                variant: 'destructive',
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="container flex items-center justify-center min-h-screen py-8">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <div className="flex justify-center mb-2">
                        <img
                            src="/lovable-uploads/bad-Photoroom.png"
                            alt="Shuttlecock Logo"
                            className="w-25 h-24 object-contain"
                        />
                    </div>
                    <CardTitle className="text-2xl font-bold text-court-orange">เข้าสู่ระบบ</CardTitle>
                    <CardDescription>กรอกข้อมูลเพื่อเข้าสู่ระบบจองคอร์ทแบดมินตัน</CardDescription>
                </CardHeader>
                <CardContent>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                            <FormField
                                control={form.control}
                                name="studentId"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>รหัสนักศึกษา</FormLabel>
                                        <FormControl>
                                            <Input placeholder="กรอกรหัสนักศึกษา" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={form.control}
                                name="password"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>รหัสผ่าน</FormLabel>
                                        <FormControl>
                                            <div className="relative">
                                                <Input
                                                    type={showPassword ? 'text' : 'password'}
                                                    placeholder="กรอกรหัสผ่าน"
                                                    {...field}
                                                />
                                                <Button
                                                    type="button"
                                                    variant="ghost"
                                                    size="icon"
                                                    className="absolute right-0 top-0 h-full px-3"
                                                    onClick={() => setShowPassword(!showPassword)}
                                                >
                                                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                                                </Button>
                                            </div>
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <Button
                                type="submit"
                                className="w-full bg-court-orange hover:bg-court-orange/90"
                                disabled={isLoading}
                            >
                                {isLoading ? 'กำลังดำเนินการ...' : 'เข้าสู่ระบบ'}
                            </Button>
                        </form>
                    </Form>
                </CardContent>
                <CardFooter className="flex justify-center">
                    <p className="text-sm text-center">
                        ยังไม่มีบัญชี?{' '}
                        <Link to="/register" className="text-court-orange hover:underline">
                            สมัครสมาชิก
                        </Link>
                    </p>
                </CardFooter>
            </Card>
        </div>
    );
};

export default Login;