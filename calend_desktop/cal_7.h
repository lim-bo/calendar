#ifndef CAL_7_H
#define CAL_7_H
#include <QHBoxLayout>
#include <QMessageBox>
#include <QWidget>
#include "event_entry.h"
#include "client.h"
#include "eventdata.h"

namespace Ui {
class cal_7;
}

class cal_7 : public QWidget
{
    Q_OBJECT

public:
    explicit cal_7(QString uid,QWidget *parent = nullptr);
    ~cal_7();
private slots:
    void showEvent(QShowEvent *event);
    void loadEventsForWeek();
    void eventDeleted(event_entry* ev);
    void eventEdited(event_entry* ev);

private:
    Ui::cal_7 *ui;
    Client cli;
    QString uid;
};

#endif // CAL_7_H
